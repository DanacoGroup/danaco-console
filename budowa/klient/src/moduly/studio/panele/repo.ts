/**
 * Panel Session Repository okna Studia — historia wersji dokumentu w
 * repozytorium sesji: autor, czas, opis zmiany, przywracanie i rozgałęzianie.
 *
 * Rodziny komend tego panelu (`studio.version.*`, `studio.repository.*`,
 * `studio.journal.*`, `studio.branch.*`) operują na `documentId`, którego
 * `ZaleznosciPanelu` nie niesie — panel poznaje go wyłącznie ze zdarzenia
 * `studio.document.changed`, którego kontrakt wprost przypisuje do odświeżania
 * tego panelu. Do pierwszego takiego zdarzenia po zamontowaniu panel czeka,
 * nie zgaduje dokumentu z innego źródła.
 *
 * Podgląd i porównanie wersji zostają poza tym plikiem: wymagają rodzin
 * `studio.diff.*` i `studio.preview.*`, przydzielonych innym zakresom prac.
 */

import {
  ChangeKind,
  Command,
  EventType,
  StudioAuthor,
  type ErrorInfo,
  type StudioVersion,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, zeZnacznika } from '../narzedzia.ts';
import type { MontazPanelu, ZamontowanyPanel, ZaleznosciPanelu } from './umowa.ts';
import { komunikatEksportu, komunikatGalezi, tresciRepo, znacznikWersji } from './repo-tresci.ts';
import { opisOdmowy, zaloguj } from '../odmowa.ts';

type StanGlowny =
  | { rodzaj: 'brakOkna' }
  | { rodzaj: 'oczekiwanieDokumentu' }
  | { rodzaj: 'wczytywanie'; dokument: string }
  | { rodzaj: 'odmowaListy'; dokument: string; blad?: ErrorInfo }
  | { rodzaj: 'lista'; dokument: string; wersje: StudioVersion[] };

type RodzajFormularza = 'rozgalez' | 'etykieta';
type StanFormularza = { akcja: RodzajFormularza; wersjaId: string } | null;
type StanKomunikatu = { rodzaj: 'sukces' | 'blad'; tekst: string } | null;
type StanAkcjiWiersza = { wersjaId: string; akcja: 'przywroc' | RodzajFormularza } | null;

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function opisAutora(autor: StudioAuthor | undefined): string {
  if (autor === StudioAuthor.Model) return tresciRepo.wiersz.autorModel;
  if (autor === StudioAuthor.Uzytkownik) return tresciRepo.wiersz.autorUzytkownik;
  return tresciRepo.wiersz.autorNieznany;
}

function formatCzas(znacznikCzasu: number): string {
  return new Date(znacznikCzasu).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

function formatBajty(bajty: number): string {
  if (bajty < 1024) return `${bajty} B`;
  const kb = bajty / 1024;
  if (kb < 1024) return `${kb.toFixed(1)} kB`;
  return `${(kb / 1024).toFixed(1)} MB`;
}

export const montujPanelRepo: MontazPanelu = (wezel, zaleznosci) => zamontuj(wezel, zaleznosci);

function zamontuj(wezel: HTMLElement, zaleznosci: ZaleznosciPanelu): ZamontowanyPanel {
  wezel.classList.add('sta-okno');
  wezel.id = 'panel-repo';
  wezel.setAttribute('role', 'tabpanel');
  wezel.setAttribute('aria-labelledby', 'karta-repo');
  wezel.setAttribute('data-nazwa', tresciRepo.panel.tytul);

  let zdjete = false;
  let stan: StanGlowny = zaleznosci.idOkna === null ? { rodzaj: 'brakOkna' } : { rodzaj: 'oczekiwanieDokumentu' };
  let formularz: StanFormularza = null;
  let komunikat: StanKomunikatu = null;
  let akcjaWiersza: StanAkcjiWiersza = null;
  let eksportTrwa = false;

  const znacznik = el('span', { klasa: 'sta-okno-znacznik', hidden: true });
  const tresc = el('div', { klasa: 'sta-okno-tresc' });
  const naglowek = el('header', { klasa: 'sta-okno-belka' }, [
    el('span', { klasa: 'sta-okno-tytul' }, [znak('historia'), el('b', { tekst: tresciRepo.panel.tytul })]),
    znacznik,
  ]);
  wezel.replaceChildren(naglowek, tresc);

  async function wczytaj(dokument: string): Promise<void> {
    stan = { rodzaj: 'wczytywanie', dokument };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioRepositoryList, { documentId: dokument });
    if (zdjete) return;
    stan =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'odmowaListy', dokument, blad: wynik.blad }
        : { rodzaj: 'lista', dokument, wersje: wynik.wynik.versions };
    odswiez();
  }

  async function przywroc(dokument: string, wersja: StudioVersion): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'przywroc' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioRepositoryRestore, {
      documentId: dokument,
      versionId: wersja.id,
    });
    if (zdjete) return;
    akcjaWiersza = null;
    if (!wynik.udany) {
      zaloguj(wynik.blad, 'repo.przywroc');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.przywroc };
      odswiez();
      return;
    }
    komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.przywrocono };
    await wczytaj(dokument);
  }

  async function zalozGalaz(dokument: string, wersja: StudioVersion, nazwa: string): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'rozgalez' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBranchCreate, {
      documentId: dokument,
      fromVersionId: wersja.id,
      name: nazwa,
    });
    if (zdjete) return;
    akcjaWiersza = null;
    formularz = null;
    zaloguj(wynik.blad, 'repo.galaz');
    komunikat =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'blad', tekst: tresciRepo.odmowa.galaz }
        : { rodzaj: 'sukces', tekst: komunikatGalezi(wynik.wynik.branch.name) };
    odswiez();
  }

  async function ustawEtykiete(wersja: StudioVersion, etykieta: string, kluczowa: boolean): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'etykieta' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioVersionLabelSet, {
      versionId: wersja.id,
      label: etykieta === '' ? undefined : etykieta,
      milestone: kluczowa,
    });
    if (zdjete) return;
    akcjaWiersza = null;
    formularz = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.etykieta');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.etykieta };
      odswiez();
      return;
    }
    if (stan.rodzaj === 'lista') {
      const podmieniona = wynik.wynik.version;
      stan = { ...stan, wersje: stan.wersje.map((w) => (w.id === podmieniona.id ? podmieniona : w)) };
    }
    komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.etykietaZapisana };
    odswiez();
  }

  async function eksportuj(dokument: string): Promise<void> {
    eksportTrwa = true;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioRepositoryExport, { documentId: dokument });
    if (zdjete) return;
    eksportTrwa = false;
    zaloguj(wynik.blad, 'repo.eksport');
    komunikat =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'blad', tekst: tresciRepo.odmowa.eksport }
        : { rodzaj: 'sukces', tekst: komunikatEksportu(wynik.wynik.entries, formatBajty(wynik.wynik.sizeBytes)) };
    odswiez();
  }

  function widokLadowanie(opis: string): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
        el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
        opis,
      ]),
    ];
  }

  function widokOdmowa(tytul: string, blad?: ErrorInfo): HTMLElement[] {
    const dzieciTresci = [el('b', { tekst: tytul }), el('span', { tekst: opisOdmowy(blad, 'repo') })];
    return [
      el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
        el('div', { klasa: 'dn-alert-tresc' }, dzieciTresci),
      ]),
    ];
  }

  function widokKomunikat(k: NonNullable<StanKomunikatu>): HTMLElement {
    const udany = k.rodzaj === 'sukces';
    return el(
      'div',
      { klasa: `dn-alert dn-alert--wstega ${udany ? 'dn-alert--sukces' : 'dn-alert--blad'}`, role: udany ? 'status' : 'alert' },
      [el('span', { klasa: 'dn-alert-tresc', tekst: k.tekst })],
    );
  }

  function widokPusto(): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty' }, [
        znak('pusto'),
        el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tresciRepo.pusto.tytul }),
        el('span', { klasa: 'dn-pusty-stan-opis', tekst: tresciRepo.pusto.opis }),
      ]),
    ];
  }

  function formularzGalaz(dokument: string, wersja: StudioVersion, busy: boolean): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      required: true,
      placeholder: tresciRepo.formularzGalaz.zastepcza,
      'aria-label': tresciRepo.formularzGalaz.etykieta,
    }) as HTMLInputElement;

    const przyciskZaloz = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'submit',
      disabled: busy,
      tekst: busy ? tresciRepo.formularzGalaz.zakladanie : tresciRepo.formularzGalaz.zaloz,
    });
    const przyciskAnuluj = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: tresciRepo.akcje.anuluj,
    });
    przyciskAnuluj.addEventListener('click', () => {
      formularz = null;
      odswiez();
    });

    const formularzWezel = el('form', { klasa: 'dn-pole' }, [
      el('label', { klasa: 'dn-pole-etykieta', tekst: tresciRepo.formularzGalaz.etykieta }),
      pole,
      el('div', { klasa: 'dn-pas-dzialan' }, [przyciskZaloz, przyciskAnuluj]),
    ]);
    formularzWezel.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      const nazwa = pole.value.trim();
      if (nazwa === '') return;
      void zalozGalaz(dokument, wersja, nazwa);
    });
    return formularzWezel;
  }

  function formularzEtykieta(wersja: StudioVersion, busy: boolean): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': tresciRepo.formularzEtykieta.etykieta,
    }) as HTMLInputElement;
    pole.value = wersja.label ?? '';

    const zaznacz = el('input', { type: 'checkbox', id: `repo-kluczowa-${wersja.id}` }) as HTMLInputElement;
    zaznacz.checked = wersja.milestone === true;
    const zaznaczEtykieta = el('label', { for: `repo-kluczowa-${wersja.id}` }, [
      zaznacz,
      ` ${tresciRepo.formularzEtykieta.kluczowa}`,
    ]);

    const przyciskZapisz = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'submit',
      disabled: busy,
      tekst: busy ? tresciRepo.formularzEtykieta.zapisywanie : tresciRepo.formularzEtykieta.zapisz,
    });
    const przyciskAnuluj = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: tresciRepo.akcje.anuluj,
    });
    przyciskAnuluj.addEventListener('click', () => {
      formularz = null;
      odswiez();
    });

    const formularzWezel = el('form', { klasa: 'dn-pole' }, [
      el('label', { klasa: 'dn-pole-etykieta', tekst: tresciRepo.formularzEtykieta.etykieta }),
      pole,
      zaznaczEtykieta,
      el('div', { klasa: 'dn-pas-dzialan' }, [przyciskZapisz, przyciskAnuluj]),
    ]);
    formularzWezel.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      void ustawEtykiete(wersja, pole.value.trim(), zaznacz.checked);
    });
    return formularzWezel;
  }

  function wiersz(dokument: string, wersja: StudioVersion, indeks: number, liczbaWersji: number): HTMLElement {
    const biezaca = indeks === 0;
    const busyTegoWiersza = akcjaWiersza?.wersjaId === wersja.id;

    const naglowekWiersza = el('div', { klasa: 'st-panel-wiersz' }, [
      el('span', {
        klasa: `dn-kropka ${biezaca ? 'dn-kropka--sukces' : 'dn-kropka--neutralna'}`,
        'aria-hidden': 'true',
      }),
      el('b', { tekst: `Wersja ${liczbaWersji - indeks}` }),
      el('span', { klasa: 'dn-meta', tekst: `${formatCzas(wersja.createdAt)} · ${opisAutora(wersja.author)}` }),
    ]);

    const dzieci: HTMLElement[] = [naglowekWiersza];

    if (biezaca) {
      dzieci.push(el('div', { klasa: 'dn-nota', tekst: tresciRepo.akcje.biezaca }));
    } else if (wersja.summary !== undefined && wersja.summary !== '') {
      dzieci.push(el('div', { klasa: 'dn-nota', tekst: `„${wersja.summary}”` }));
    }

    if (wersja.milestone === true || (wersja.label !== undefined && wersja.label !== '')) {
      const tekstPlakietki = wersja.milestone === true ? `★ ${wersja.label ?? tresciRepo.wiersz.kluczowa}` : (wersja.label as string);
      dzieci.push(
        el('div', {}, [
          el('span', { klasa: `dn-plakietka${wersja.milestone === true ? ' dn-plakietka--sukces' : ''}`, tekst: tekstPlakietki }),
        ]),
      );
    }

    const przyciskPrzywroc = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: biezaca || busyTegoWiersza,
      tekst: akcjaWiersza?.wersjaId === wersja.id && akcjaWiersza.akcja === 'przywroc'
        ? tresciRepo.akcje.przywracanie
        : tresciRepo.akcje.przywroc,
    });
    przyciskPrzywroc.addEventListener('click', () => void przywroc(dokument, wersja));

    const przyciskRozgalez = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: busyTegoWiersza,
      tekst: tresciRepo.akcje.rozgalez,
    });
    przyciskRozgalez.addEventListener('click', () => {
      formularz =
        formularz?.wersjaId === wersja.id && formularz.akcja === 'rozgalez'
          ? null
          : { akcja: 'rozgalez', wersjaId: wersja.id };
      odswiez();
    });

    const przyciskWiecej = el('button', {
      klasa: 'dn-btn-ikona dn-btn-ikona--sm',
      type: 'button',
      'aria-label': tresciRepo.akcje.wiecej,
      disabled: busyTegoWiersza,
      tekst: '⋯',
    });
    przyciskWiecej.addEventListener('click', () => {
      formularz =
        formularz?.wersjaId === wersja.id && formularz.akcja === 'etykieta'
          ? null
          : { akcja: 'etykieta', wersjaId: wersja.id };
      odswiez();
    });

    dzieci.push(el('div', { klasa: 'dn-wersja-akcje' }, [przyciskPrzywroc, przyciskRozgalez, przyciskWiecej]));

    if (formularz !== null && formularz.wersjaId === wersja.id) {
      dzieci.push(
        formularz.akcja === 'rozgalez'
          ? formularzGalaz(dokument, wersja, busyTegoWiersza)
          : formularzEtykieta(wersja, busyTegoWiersza),
      );
    }

    return el('div', { klasa: 'dn-wersja' }, dzieci);
  }

  function stopkaEksportu(dokument: string): HTMLElement {
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm st-pole-rozciagniete',
      type: 'button',
      disabled: eksportTrwa,
      tekst: eksportTrwa ? tresciRepo.stopka.eksportowanie : tresciRepo.stopka.eksportuj,
    });
    przycisk.addEventListener('click', () => void eksportuj(dokument));
    return el('div', { klasa: 'dn-wersja-stopka' }, [przycisk]);
  }

  function widokLista(dokument: string, wersje: StudioVersion[]): HTMLElement[] {
    const elementy: HTMLElement[] = [];
    if (komunikat !== null) elementy.push(widokKomunikat(komunikat));
    if (wersje.length === 0) {
      elementy.push(...widokPusto());
      return elementy;
    }
    wersje.forEach((wersja, indeks) => elementy.push(wiersz(dokument, wersja, indeks, wersje.length)));
    elementy.push(stopkaEksportu(dokument));
    return elementy;
  }

  function zawartosc(): HTMLElement[] {
    switch (stan.rodzaj) {
      case 'brakOkna':
        return widokOdmowa(tresciRepo.stany.brakOkna);
      case 'oczekiwanieDokumentu':
        return widokLadowanie(tresciRepo.stany.oczekiwanieDokumentu);
      case 'wczytywanie':
        return widokLadowanie(tresciRepo.stany.wczytywanie);
      case 'odmowaListy':
        return widokOdmowa(tresciRepo.odmowa.lista, stan.blad);
      case 'lista':
        return widokLista(stan.dokument, stan.wersje);
    }
  }

  function odswiez(): void {
    tresc.replaceChildren(...zawartosc());
    znacznik.hidden = stan.rodzaj !== 'lista';
    znacznik.textContent = stan.rodzaj === 'lista' ? znacznikWersji(stan.wersje.length) : '';
  }

  odswiez();

  const odsubskrybuj: (() => void) | null =
    zaleznosci.idOkna === null
      ? null
      : zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
          if (zdjete || zdarzenie.document.windowId !== zaleznosci.idOkna) return;
          if (zdarzenie.change === ChangeKind.Deleted) {
            stan = { rodzaj: 'oczekiwanieDokumentu' };
            komunikat = null;
            odswiez();
            return;
          }
          void wczytaj(zdarzenie.document.id);
        });

  return {
    zdejmij() {
      zdjete = true;
      odsubskrybuj?.();
      wezel.replaceChildren();
    },
  };
}
