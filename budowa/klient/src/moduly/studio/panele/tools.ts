/**
 * Panel Tools Panel okna Studia. Kształt wzięty ze źródła prawdy
 * `design/05-okna/moduly/studio.html`, `section#panel-tools`: przełącznik
 * zakresu, etykieta rozmiaru zakresu, wykaz operacji rozpisany kategoriami
 * i przycisk „Uruchom operację”.
 *
 * Wykaz operacji przychodzi komendą `studio.operation.list` — żadna pozycja
 * nie jest wpisana w kod. Wybraną pozycję puszcza `studio.contextual.op`,
 * własną zakłada `studio.operation.save`, a zdejmuje `studio.operation.delete`.
 * Osobnym torem stoi `studio.search.semantic`: szuka w dokumencie fragmentów
 * bliskich zapytaniu, gdy dosłowne dopasowanie nie wystarcza.
 *
 * Dokument, na którym operacje działają, panel poznaje ze zdarzenia
 * `studio.document.changed` ograniczonego do własnego okna — kontrakt nie
 * niesie komendy zwracającej wprost dokument otwarty w oknie, a to samo
 * zdarzenie dochodzi do Studio Editor przy każdym założeniu i zapisie.
 */

import {
  ChangeKind,
  Command,
  EventType,
  StudioOperationScope,
  type ErrorInfo,
  type StudioDocument,
  type StudioOperation,
  type StudioSemanticMatch,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony } from '../ikony.ts';
import { el, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { tresciNarzedzi as T } from './tools-tresci.ts';
import type { MontazPanelu } from './umowa.ts';
import { opisOdmowy } from '../odmowa.ts';
import { zLiczba } from '../liczebnik.ts';

type StanBiegu =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'uruchamia' }
  | { rodzaj: 'wynik'; tresc?: string; propozycja?: string }
  | { rodzaj: 'odmowa'; etykieta: string; blad?: ErrorInfo };

type StanUsuwania =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'usuwa' }
  | { rodzaj: 'fabryczna' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

type StanZapisu =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'zapisuje' }
  | { rodzaj: 'zapisana' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

type StanSzukania =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'szuka' }
  | { rodzaj: 'wynik'; dopasowania: StudioSemanticMatch[] }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo };

type Stan =
  | { rodzaj: 'brakOkna' }
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'odmowaListy'; blad?: ErrorInfo }
  | { rodzaj: 'gotowy' };

function znak(rysunek: keyof typeof ikony): SVGElement {
  const wezelZnaku = zeZnacznika(ikony[rysunek]);
  wezelZnaku.setAttribute('aria-hidden', 'true');
  return wezelZnaku;
}

function alert(etykieta: string, blad?: ErrorInfo): HTMLElement {
  return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
    el('span', { klasa: 'dn-alert-tresc' }, [
      el('b', { tekst: etykieta }),
      el('span', { tekst: opisOdmowy(blad, 'tools') }),
    ]),
  ]);
}

function wierszPulsu(etykieta: string): HTMLElement {
  return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
    el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
    etykieta,
  ]);
}

export const montujPanelTools: MontazPanelu = (wezel, zaleznosci) => {
  let zdjete = false;
  let stan: Stan = zaleznosci.idOkna === null ? { rodzaj: 'brakOkna' } : { rodzaj: 'ladowanie' };

  let operacje: StudioOperation[] = [];
  let dokument: StudioDocument | null = null;
  let zakres: StudioOperationScope = StudioOperationScope.Selection;
  let wybrana: string | null = null;
  let bieg: StanBiegu = { rodzaj: 'spoczynek' };
  let usuwanie: StanUsuwania = { rodzaj: 'spoczynek' };
  let zapis: StanZapisu = { rodzaj: 'spoczynek' };
  let szukanie: StanSzukania = { rodzaj: 'spoczynek' };

  /* Treść pól trzymana poza węzłami: widok przebudowuje się w całości po każdej
     odpowiedzi, a wpis Operatora ma to przetrwać. */
  let nazwaWlasna = '';
  let kategoriaWlasna = '';
  let trescWlasna = '';
  let zapytanie = '';

  /* Uchwyty do przycisków zależnych od wpisu — trzymane, żeby wpis zmieniał
     samą ich dostępność, bez przebudowy widoku. */
  let przyciskZapisu: HTMLButtonElement | null = null;
  let przyciskSzukania: HTMLButtonElement | null = null;

  const tresc = el('div', { klasa: 'sta-okno-tresc st-panel-lista' });
  wezel.classList.add('sta-okno');
  wezel.replaceChildren(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('klucz'), el('b', { tekst: T.panel.tytul })]),
    ]),
    tresc,
  );

  function operacjaWybrana(): StudioOperation | null {
    return operacje.find((op) => op.id === wybrana) ?? null;
  }

  function przelacznikZakresu(): HTMLElement {
    function przycisk(wartosc: StudioOperationScope, etykieta: string): HTMLElement {
      const aktywny = zakres === wartosc;
      const guzik = el('button', {
        klasa: aktywny ? 'dn-btn dn-btn--zarys dn-btn--sm' : 'dn-btn dn-btn--duch dn-btn--sm',
        type: 'button',
        'aria-pressed': aktywny ? 'true' : 'false',
        tekst: etykieta,
      });
      guzik.addEventListener('click', () => {
        if (zakres === wartosc) return;
        zakres = wartosc;
        odswiezGotowy();
      });
      return guzik;
    }
    return el(
      'div',
      { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': T.zakres.etykietaPrzelacznika },
      [przycisk(StudioOperationScope.Selection, T.zakres.zaznaczenie), przycisk(StudioOperationScope.Document, T.zakres.dokument)],
    );
  }

  function etykietaRozmiaru(): HTMLElement {
    let opis: string;
    if (dokument === null) {
      opis = T.zakres.brakDokumentu;
    } else if (zakres === StudioOperationScope.Selection) {
      opis = T.zakres.brakZaznaczenia;
    } else {
      const trescDokumentu = dokument.content ?? '';
      const znaki = trescDokumentu.length;
      const slowa = trescDokumentu.split(/\s+/).filter((czlon) => czlon.length > 0).length;
      opis = `${T.zakres.etykietaRozmiaru} ${zLiczba(znaki, T.zakres.jednostkaZnaki)} · ${zLiczba(slowa, T.zakres.jednostkaSlowa)}`;
    }
    return el('div', { klasa: 'pt-etykieta', tekst: opis });
  }

  function wierszOperacji(op: StudioOperation): HTMLElement {
    const aktywna = wybrana === op.id;
    const dzieci: Dziecko[] = [
      op.name,
      op.builtin ? null : el('span', { klasa: 'dn-meta', tekst: T.operacje.wlasna }),
      el('span', { klasa: 'dn-meta', tekst: aktywna ? T.operacje.znakWyboru : T.operacje.znakWiersza }),
    ];
    /* `<div>`, jak w źródle kształtu — `.st-panel-wiersz` nie niesie resetu
       wyglądu natywnego przycisku, więc rolę i klawiaturę dokłada się wprost
       zamiast przez `<button>`, który wniósłby obramowanie przeglądarki. */
    const wiersz = el(
      'div',
      { klasa: 'st-panel-wiersz', role: 'button', tabindex: '0', 'aria-pressed': aktywna ? 'true' : 'false' },
      dzieci,
    );
    const przelacz = (): void => {
      wybrana = aktywna ? null : op.id;
      usuwanie = { rodzaj: 'spoczynek' };
      odswiezGotowy();
    };
    wiersz.addEventListener('click', przelacz);
    wiersz.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key !== 'Enter' && zdarzenie.key !== ' ') return;
      zdarzenie.preventDefault();
      przelacz();
    });
    return wiersz;
  }

  function listaOperacji(): HTMLElement[] {
    if (operacje.length === 0) {
      return [el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: T.operacje.brakWykazu })];
    }
    const wezly: HTMLElement[] = [];
    let ostatniaKategoria: string | null = null;
    for (const op of operacje) {
      if (op.category !== ostatniaKategoria) {
        wezly.push(el('div', { klasa: 'pt-etykieta', tekst: op.category }));
        ostatniaKategoria = op.category;
      }
      wezly.push(wierszOperacji(op));
    }
    return wezly;
  }

  function widokBiegu(): HTMLElement[] {
    switch (bieg.rodzaj) {
      case 'spoczynek':
        return [];
      case 'uruchamia':
        return [wierszPulsu(T.uruchom.wBiegu)];
      case 'wynik': {
        const wezly = [el('div', { klasa: 'dn-nota', tekst: bieg.tresc ?? T.uruchom.brakWyniku })];
        /* Propozycję zmiany ogląda się w karcie Diff/Grep Panel — tutaj zostaje
           samo zdanie o tym, że powstała, żeby wynik nie wyglądał na zgubiony. */
        if (bieg.propozycja !== undefined) {
          wezly.push(el('div', { klasa: 'dn-nota', tekst: T.uruchom.propozycja }));
        }
        return wezly;
      }
      case 'odmowa':
        return [alert(bieg.etykieta, bieg.blad)];
    }
  }

  function przyciskUruchom(): HTMLElement {
    const zablokowany = dokument === null || wybrana === null || bieg.rodzaj === 'uruchamia';
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--atrament st-odsun-sekcja',
      type: 'button',
      disabled: zablokowany,
      tekst: bieg.rodzaj === 'uruchamia' ? T.uruchom.wBiegu : T.uruchom.przycisk,
    });
    przycisk.addEventListener('click', () => void uruchomOperacje());
    return przycisk;
  }

  /* Zdejmowanie stoi osobno, nie w wierszu wykazu: przycisk w wierszu z rolą
     przycisku dawałby dwa sterowania jedno w drugim i klawiatura gubiłaby, na
     którym z nich stoi. */
  function widokUsuwania(): HTMLElement[] {
    const op = operacjaWybrana();
    if (op === null || op.builtin) return [];
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      disabled: usuwanie.rodzaj === 'usuwa',
      tekst: usuwanie.rodzaj === 'usuwa' ? T.usun.wBiegu : T.usun.przycisk,
    });
    przycisk.addEventListener('click', () => void usunOperacje(op.id));
    const wezly = [przycisk];
    if (usuwanie.rodzaj === 'fabryczna') wezly.push(el('div', { klasa: 'dn-nota', tekst: T.usun.fabryczna }));
    if (usuwanie.rodzaj === 'odmowa') wezly.push(alert(T.odmowa.usuniecie, usuwanie.blad));
    return wezly;
  }

  function polePisane(
    etykieta: string,
    zastepcza: string,
    wartosc: string,
    zapisz: (nowa: string) => void,
  ): HTMLInputElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': etykieta,
      placeholder: zastepcza,
      value: wartosc,
    }) as HTMLInputElement;
    pole.addEventListener('input', () => {
      zapisz(pole.value);
      odswiezStanPrzyciskow();
    });
    return pole;
  }

  function sekcjaZapisu(): HTMLElement[] {
    const naglowek = el('div', { klasa: 'pt-etykieta st-odsun-sekcja', tekst: T.zapis.naglowek });
    const poleNazwy = polePisane(T.zapis.etykietaNazwy, T.zapis.zastepczaNazwa, nazwaWlasna, (nowa) => {
      nazwaWlasna = nowa;
    });
    const poleKategorii = polePisane(T.zapis.etykietaKategorii, T.zapis.zastepczaKategoria, kategoriaWlasna, (nowa) => {
      kategoriaWlasna = nowa;
    });
    const poleTresci = polePisane(T.zapis.etykietaTresci, T.zapis.zastepczaTresc, trescWlasna, (nowa) => {
      trescWlasna = nowa;
    });
    przyciskZapisu = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: !zapisGotowy(),
      tekst: zapis.rodzaj === 'zapisuje' ? T.zapis.wBiegu : T.zapis.przycisk,
    }) as HTMLButtonElement;
    przyciskZapisu.addEventListener('click', () => void zapiszOperacje());

    const wezly = [
      naglowek,
      el('div', { klasa: 'st-panel-wiersz' }, [poleNazwy]),
      el('div', { klasa: 'st-panel-wiersz' }, [poleKategorii]),
      el('div', { klasa: 'st-panel-wiersz' }, [el('div', { klasa: 'dn-pole-zestaw' }, [poleTresci, przyciskZapisu])]),
    ];
    if (zapis.rodzaj === 'zapisana') wezly.push(el('div', { klasa: 'dn-nota', tekst: T.zapis.zapisana }));
    if (zapis.rodzaj === 'odmowa') wezly.push(alert(T.odmowa.zapis, zapis.blad));
    return wezly;
  }

  function opisDopasowania(dopasowanie: StudioSemanticMatch): string {
    const bliskosc = `${T.szukanie.bliskosc} ${dopasowanie.score.toFixed(2)}`;
    if (dopasowanie.line === undefined) return bliskosc;
    return `${T.szukanie.wiersz} ${dopasowanie.line} · ${bliskosc}`;
  }

  function widokSzukania(): HTMLElement[] {
    switch (szukanie.rodzaj) {
      case 'spoczynek':
        return [];
      case 'szuka':
        return [wierszPulsu(T.szukanie.wBiegu)];
      case 'wynik':
        if (szukanie.dopasowania.length === 0) {
          return [el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: T.szukanie.brakDopasowan })];
        }
        return szukanie.dopasowania.map((dopasowanie) =>
          el('div', { klasa: 'st-panel-wiersz' }, [
            el('span', { tekst: dopasowanie.text }),
            el('span', { klasa: 'dn-meta', tekst: opisDopasowania(dopasowanie) }),
          ]),
        );
      case 'odmowa':
        return [alert(T.odmowa.szukanie, szukanie.blad)];
    }
  }

  function sekcjaSzukania(): HTMLElement[] {
    const pole = el('input', {
      type: 'search',
      'aria-label': T.szukanie.etykieta,
      placeholder: T.szukanie.zastepczaTresc,
      value: zapytanie,
      disabled: dokument === null,
    }) as HTMLInputElement;
    pole.addEventListener('input', () => {
      zapytanie = pole.value;
      odswiezStanPrzyciskow();
    });
    pole.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Enter') void szukajZnaczeniowo();
    });
    przyciskSzukania = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      disabled: !szukanieGotowe(),
      tekst: szukanie.rodzaj === 'szuka' ? T.szukanie.wBiegu : T.szukanie.przycisk,
    }) as HTMLButtonElement;
    przyciskSzukania.addEventListener('click', () => void szukajZnaczeniowo());

    return [
      el('div', { klasa: 'pt-etykieta st-odsun-sekcja', tekst: T.szukanie.naglowek }),
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('label', { klasa: 'dn-szukaj st-szukaj-odsun' }, [znak('szukaj'), pole]),
        przyciskSzukania,
      ]),
      ...widokSzukania(),
    ];
  }

  function zapisGotowy(): boolean {
    return (
      zapis.rodzaj !== 'zapisuje' &&
      nazwaWlasna.trim() !== '' &&
      kategoriaWlasna.trim() !== '' &&
      trescWlasna.trim() !== ''
    );
  }

  function szukanieGotowe(): boolean {
    return szukanie.rodzaj !== 'szuka' && dokument !== null && zapytanie.trim() !== '';
  }

  /* Przebudowa całego widoku przy każdym znaku gubiłaby ognisko pola, więc wpis
     rusza wyłącznie dostępnością przycisków, do których się odnosi. */
  function odswiezStanPrzyciskow(): void {
    if (przyciskZapisu !== null) przyciskZapisu.disabled = !zapisGotowy();
    if (przyciskSzukania !== null) przyciskSzukania.disabled = !szukanieGotowe();
  }

  function widokGotowy(): HTMLElement[] {
    return [
      przelacznikZakresu(),
      etykietaRozmiaru(),
      ...listaOperacji(),
      /* Prototyp pokazuje przy części operacji wybór wariantu („Zmień styl ▾”).
         Operacje przychodzą bez opisu swoich ustawień, więc zamiast pustego
         rozwinięcia stoi zdanie nazywające tę niegotowość. */
      el('div', { klasa: 'dn-nota', tekst: T.operacje.warianty }),
      ...widokBiegu(),
      przyciskUruchom(),
      ...widokUsuwania(),
      ...sekcjaSzukania(),
      ...sekcjaZapisu(),
    ];
  }

  function zawartosc(s: Stan): HTMLElement[] {
    switch (s.rodzaj) {
      case 'brakOkna':
        return [alert(T.odmowa.brakOkna)];
      case 'ladowanie':
        return [wierszPulsu(T.ladowanie)];
      case 'odmowaListy':
        return [alert(T.odmowa.lista, s.blad)];
      case 'gotowy':
        return widokGotowy();
    }
  }

  function odswiez(s: Stan): void {
    stan = s;
    przyciskZapisu = null;
    przyciskSzukania = null;
    tresc.replaceChildren(...zawartosc(s));
  }

  function odswiezGotowy(): void {
    if (stan.rodzaj !== 'gotowy') return;
    odswiez(stan);
  }

  async function zaladujOperacje(pierwsze: boolean): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioOperationList, {});
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      /* Odmowa przy odświeżeniu po zapisie nie zwija całej karty — wykaz sprzed
         wywołania zostaje, a zdanie o odmowie idzie tam, gdzie czynność stała. */
      if (pierwsze) odswiez({ rodzaj: 'odmowaListy', blad: wynik.blad });
      return;
    }
    operacje = wynik.wynik.operations;
    if (wybrana !== null && operacjaWybrana() === null) wybrana = null;
    odswiez({ rodzaj: 'gotowy' });
  }

  async function uruchomOperacje(): Promise<void> {
    const idOkna = zaleznosci.idOkna;
    const dokumentBiezacy = dokument;
    const operacjaId = wybrana;
    if (idOkna === null || dokumentBiezacy === null || operacjaId === null) return;

    bieg = { rodzaj: 'uruchamia' };
    odswiezGotowy();

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioContextualOp, {
      windowId: idOkna,
      documentId: dokumentBiezacy.id,
      actionId: operacjaId,
      scope: zakres,
    });
    if (zdjete) return;
    bieg = wynik.udany
      ? { rodzaj: 'wynik', tresc: wynik.wynik?.resultText, propozycja: wynik.wynik?.proposalId }
      : { rodzaj: 'odmowa', etykieta: T.odmowa.uruchomienie, blad: wynik.blad };
    odswiezGotowy();
  }

  async function zapiszOperacje(): Promise<void> {
    if (!zapisGotowy()) return;
    zapis = { rodzaj: 'zapisuje' };
    odswiezGotowy();

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioOperationSave, {
      name: nazwaWlasna.trim(),
      category: kategoriaWlasna.trim(),
      prompt: trescWlasna.trim(),
    });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      zapis = { rodzaj: 'odmowa', blad: wynik.blad };
      odswiezGotowy();
      return;
    }
    nazwaWlasna = '';
    kategoriaWlasna = '';
    trescWlasna = '';
    zapis = { rodzaj: 'zapisana' };
    await zaladujOperacje(false);
  }

  async function usunOperacje(operacjaId: string): Promise<void> {
    usuwanie = { rodzaj: 'usuwa' };
    odswiezGotowy();

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioOperationDelete, { operationId: operacjaId });
    if (zdjete) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      usuwanie = { rodzaj: 'odmowa', blad: wynik.blad };
      odswiezGotowy();
      return;
    }
    if (!wynik.wynik.deleted) {
      usuwanie = { rodzaj: 'fabryczna' };
      odswiezGotowy();
      return;
    }
    usuwanie = { rodzaj: 'spoczynek' };
    wybrana = null;
    await zaladujOperacje(false);
  }

  async function szukajZnaczeniowo(): Promise<void> {
    const dokumentBiezacy = dokument;
    if (!szukanieGotowe() || dokumentBiezacy === null) return;
    szukanie = { rodzaj: 'szuka' };
    odswiezGotowy();

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioSearchSemantic, {
      documentId: dokumentBiezacy.id,
      query: zapytanie.trim(),
    });
    if (zdjete) return;
    szukanie =
      wynik.udany && wynik.wynik !== undefined
        ? { rodzaj: 'wynik', dopasowania: wynik.wynik.matches }
        : { rodzaj: 'odmowa', blad: wynik.blad };
    odswiezGotowy();
  }

  odswiez(stan);

  const odsubskrybujDokument = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
    if (zdjete || zaleznosci.idOkna === null || zdarzenie.document.windowId !== zaleznosci.idOkna) return;
    if (zdarzenie.change === ChangeKind.Deleted) {
      if (dokument !== null && dokument.id === zdarzenie.document.id) dokument = null;
    } else {
      dokument = zdarzenie.document;
    }
    odswiezGotowy();
  });

  if (zaleznosci.idOkna !== null) void zaladujOperacje(true);

  return {
    zdejmij() {
      zdjete = true;
      odsubskrybujDokument();
      wezel.replaceChildren();
    },
  };
};
