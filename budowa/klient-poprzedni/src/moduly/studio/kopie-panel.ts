import './dziennik-kontrola.css';

import {
  StudioBackupReason,
  StudioVersionSeries,
  type StudioAutosaveSettings,
  type StudioDocumentBackup,
  type StudioVersion,
} from '../../../../shared/contract';
import { poleLogiczne, poleTekstowe, poleWyboru } from '../../modele/kontrolki-formularza';
import {
  KOPIE_PRZED_ZAPISEM,
  kopieNazwaSzeregu,
  kopieOpiszKopie,
  kopieOpiszNastawy,
  kopieOpiszStanZapisu,
  kopieOpiszSzeregi,
  kopieStanZapisu,
  kopieZgloszeniePoZamknieciu,
  type StanZapisu,
} from './kopie-zapasu';
import type { ZrodloKontroliStudio } from './zrodlo-kontroli-studio';

/**
 * Panel autozapisu, kopii zapasowych i szeregów wersji łączy wymagania
 * trwałości dokumentu.
 */

/**
 * Czym panel pyta okno o dokument i co mu oddaje: identyfikator dokumentu
 * i okna, treść, postać oraz stan zmian niezapisanych.
 */
export interface KontekstKopii {
  idDokumentu(): string;
  /** Okno osadzenia — kopie i nastawy bywają wiązane z oknem, nie z dokumentem. */
  idOkna(): string;
  /** Treść dokumentu do odłożenia w kopii i w zapisie samoczynnym. */
  tresc(): string;
  /** Postać dokumentu w kształcie `StudioDocumentForm`; `null`, gdy okno jej nie ma. */
  postac(): unknown;
  /** Czy od ostatniego zapisu były zmiany — tego rdzeń nie widzi. */
  zmianyNiezapisane(): boolean;
  /** Przyjmuje treść i postać oddaną po przywróceniu. */
  naSkutek(tresc: string | undefined, postac: unknown): void;
}

export interface KopiePanel {
  element: HTMLElement;
  /** Odczytuje nastawy autozapisu, kopie i szeregi wersji. */
  odswiez(): Promise<void>;
  /** Wykonuje zapis samoczynny wraz z podaniem, co go wywołało. */
  zapiszSamoczynnie(powod: StudioBackupReason): Promise<boolean>;
  /** Zakłada kopię przed czynnością nieodwracalną; oddaje jej kod albo `null`. */
  kopiaPrzedCzynnoscia(): Promise<string | null>;
  /** Stan zapisu do paska statusu okna. */
  stanZapisu(): StanZapisu;
  /** Zdanie o stanie zapisu do paska statusu okna. */
  zdanieStanuZapisu(): string;
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

export function utworzKopiePanel(
  zrodlo: ZrodloKontroliStudio,
  kontekst: KontekstKopii,
): KopiePanel {
  let otwarty = false;
  let nastawy: StudioAutosaveSettings | null = null;
  let zapisWToku = false;
  let kopie: readonly StudioDocumentBackup[] = [];

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-pole-opis ms-kontrola__odpowiedz';
  odpowiedz.dataset['czynnosc'] = 'odpowiedz';

  function powiedz(tresc: string, udana: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.dataset['udana'] = udana ? 'tak' : 'nie';
  }

  function dokument(): string {
    const kod = kontekst.idDokumentu();
    if (kod === '') {
      powiedz(
        'Nie ma dokumentu czynnego. Autozapis, kopie zapasowe i szeregi wersji dotyczą jednego ' +
          'dokumentu — otwórz go albo załóż nowy.',
        false,
      );
    }
    return kod;
  }

  /* ── Wskaźnik stanu zapisu ──────────────────────────────────────────────── */

  const wskaznik = document.createElement('p');
  wskaznik.className = 'ms-kontrola__stan';
  wskaznik.dataset['czynnosc'] = 'stan-zapisu';
  wskaznik.setAttribute('role', 'status');

  function przerysujWskaznik(): void {
    const stan = kopieStanZapisu(nastawy, kontekst.zmianyNiezapisane(), zapisWToku);
    wskaznik.dataset['stan'] = stan;
    wskaznik.textContent = kopieOpiszStanZapisu(stan, nastawy);
  }

  /* ── Autozapis ──────────────────────────────────────────────────────────── */

  const czynny = poleLogiczne({
    etykieta: 'Autozapis czynny',
    opis:
      'Jawne, odwracalne ustawienie Operatora. Zapis samoczynny obejmuje treść I POSTAĆ — postać ' +
      'zgubiona przy zapisie samoczynnym byłaby gorsza niż brak autozapisu.',
  });
  const odstep = poleTekstowe({
    etykieta: 'Odstęp zapisu samoczynnego w sekundach',
    podpowiedz: 'na przykład 120',
    opis: 'Puste albo zero znaczy „bez odstępu" — wtedy zapis idzie wyłącznie przy zdarzeniach okna.',
  });
  const przyOdejsciu = poleLogiczne({ etykieta: 'Zapis przy odejściu od okna' });
  const przyZamknieciu = poleLogiczne({ etykieta: 'Zapis przy zamknięciu dokumentu' });
  const przyPrzelaczeniu = poleLogiczne({ etykieta: 'Zapis przy przełączeniu dokumentu' });
  const ileKopii = poleTekstowe({
    etykieta: 'Ile kopii zapasowych zachować',
    podpowiedz: 'na przykład 10',
  });
  const ileGodzin = poleTekstowe({
    etykieta: 'Po ilu godzinach kopia wygasa',
    podpowiedz: 'na przykład 72',
    opis: 'Zasada wygasania kopii jest ustawieniem Operatora, nie liczbą zaszytą w kodzie.',
  });

  const zapiszNastawy = przyciskPanelu('Zapisz nastawy autozapisu', 'autozapis-ustaw');
  zapiszNastawy.addEventListener('click', () => {
    void ustaw();
  });

  const zapiszTeraz = przyciskPanelu('Zapisz samoczynnie TERAZ', 'autozapis-wykonaj');
  zapiszTeraz.title =
    'Wykonuje zapis samoczynny wraz z postacią i odkłada wersję w OSOBNYM szeregu autozapisu. ' +
    'Nieudany zapis wraca nazwany, nie przemilczany.';
  zapiszTeraz.addEventListener('click', () => {
    void zapiszSamoczynnie(StudioBackupReason.Manual);
  });

  const zdanieNastaw = document.createElement('p');
  zdanieNastaw.className = 'dn-pole-opis';

  /* ── Kopie zapasowe ─────────────────────────────────────────────────────── */

  const zdanieKopii = document.createElement('p');
  zdanieKopii.className = 'dn-pole-opis ms-kontrola__zasada';
  zdanieKopii.textContent = KOPIE_PRZED_ZAPISEM;

  const zgloszenie = document.createElement('p');
  zgloszenie.className = 'dn-pole-opis';
  zgloszenie.dataset['czynnosc'] = 'zgloszenie-po-zamknieciu';

  const tylkoNiezapisane = poleLogiczne({
    etykieta: 'Tylko kopie niosące zmiany niezapisane',
    opis: 'To są kopie, po których wraca się po nagłym zamknięciu.',
  });
  tylkoNiezapisane.kontrolka.addEventListener('change', () => {
    void odswiezKopie();
  });

  const doNowego = poleLogiczne({
    etykieta: 'Przywracaj DO NOWEGO DOKUMENTU',
    opis:
      'Przywrócenie do nowego dokumentu nie rusza tego, co stoi w oknie — przywracanie samo nie ' +
      'kasuje wtedy pracy bieżącej. Zdjęcie tego znaczy przywrócenie NA MIEJSCE.',
  });
  doNowego.kontrolka.checked = true;

  const zalozKopie = przyciskPanelu('Załóż kopię zapasową teraz', 'kopia-zaloz');
  zalozKopie.addEventListener('click', () => {
    void zalozKopieRecznie();
  });

  const podsumowanieKopii = document.createElement('p');
  podsumowanieKopii.className = 'dn-pole-opis';

  const wykazKopii = document.createElement('ul');
  wykazKopii.className = 'ms-kontrola__wykaz';

  /* ── Szeregi wersji i powrót do wersji założycielskiej ──────────────────── */

  const szereg = poleWyboru(
    {
      etykieta: 'Szereg wersji',
      opis:
        'Zapisy samoczynne idą osobnym szeregiem i są w wykazie odróżnialne. Rozróżnienie ' +
        'pochodzi z rdzenia, nie z domysłu okna po braku etykiety.',
    },
    [
      { wartosc: '', etykieta: 'Oba szeregi razem' },
      { wartosc: StudioVersionSeries.Operator, etykieta: 'Tylko wersje Operatora' },
      { wartosc: StudioVersionSeries.Autosave, etykieta: 'Tylko zapisy samoczynne' },
    ],
  );
  szereg.kontrolka.addEventListener('change', () => {
    void odswiezSzeregi();
  });

  const tylkoKluczowe = poleLogiczne({ etykieta: 'Tylko wersje kluczowe' });
  tylkoKluczowe.kontrolka.addEventListener('change', () => {
    void odswiezSzeregi();
  });

  const powrot = przyciskPanelu('Wróć do wersji założycielskiej', 'powrot-zalozycielska');
  powrot.title =
    'Jedno polecenie, bez szukania wersji w wykazie. Wersje nowsze ZOSTAJĄ w historii, więc sam ' +
    'powrót da się cofnąć.';
  powrot.addEventListener('click', () => {
    void wrocDoZalozycielskiej();
  });

  const podsumowanieSzeregow = document.createElement('p');
  podsumowanieSzeregow.className = 'dn-pole-opis';

  const wykazWersji = document.createElement('ul');
  wykazWersji.className = 'ms-kontrola__wykaz';

  /* ── Złożenie panelu ────────────────────────────────────────────────────── */

  const element = document.createElement('section');
  element.className = 'ms-kontrola';
  element.dataset['panel'] = 'kopie-i-autozapis';
  element.hidden = true;
  element.setAttribute('aria-label', 'Autozapis, kopie zapasowe i szeregi wersji dokumentu');
  element.append(
    czescPanelu('Stan zapisu', [wskaznik]),
    czescPanelu('Autozapis — jawne ustawienie Operatora', [
      czynny.element,
      odstep.element,
      przyOdejsciu.element,
      przyZamknieciu.element,
      przyPrzelaczeniu.element,
      ileKopii.element,
      ileGodzin.element,
      pasPrzyciskow([zapiszNastawy, zapiszTeraz]),
      zdanieNastaw,
    ]),
    czescPanelu('Kopie zapasowe i przywrócenie', [
      zdanieKopii,
      zgloszenie,
      tylkoNiezapisane.element,
      doNowego.element,
      pasPrzyciskow([zalozKopie]),
      podsumowanieKopii,
      wykazKopii,
    ]),
    czescPanelu('Szeregi wersji i powrót do stanu pierwotnego', [
      szereg.element,
      tylkoKluczowe.element,
      pasPrzyciskow([powrot]),
      podsumowanieSzeregow,
      wykazWersji,
    ]),
    odpowiedz,
  );

  /* ── Czynności ──────────────────────────────────────────────────────────── */

  function liczba(pole: HTMLInputElement): number | undefined {
    const wpisane = pole.value.trim();
    if (wpisane === '') return undefined;
    const wartosc = Number(wpisane);
    return Number.isFinite(wartosc) && wartosc >= 0 ? wartosc : undefined;
  }

  async function odswiezNastawy(): Promise<void> {
    const wynik = await zrodlo.kontrolaAutozapisNastawy(
      kontekst.idDokumentu(),
      kontekst.idOkna(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      zdanieNastaw.textContent =
        `Nastaw autozapisu nie udało się odczytać: ${wynik.blad?.message ?? ''}`;
      return;
    }
    nastawy = wynik.wynik.settings;
    czynny.kontrolka.checked = nastawy.enabled;
    if (nastawy.intervalSeconds !== undefined) odstep.kontrolka.value = String(nastawy.intervalSeconds);
    przyOdejsciu.kontrolka.checked = nastawy.onBlur === true;
    przyZamknieciu.kontrolka.checked = nastawy.onClose === true;
    przyPrzelaczeniu.kontrolka.checked = nastawy.onSwitch === true;
    if (nastawy.backupRetentionCount !== undefined) {
      ileKopii.kontrolka.value = String(nastawy.backupRetentionCount);
    }
    if (nastawy.backupRetentionHours !== undefined) {
      ileGodzin.kontrolka.value = String(nastawy.backupRetentionHours);
    }
    zdanieNastaw.textContent = kopieOpiszNastawy(nastawy);
    przerysujWskaznik();
  }

  async function ustaw(): Promise<void> {
    const wynik = await zrodlo.kontrolaAutozapisUstaw(kontekst.idDokumentu(), kontekst.idOkna(), {
      czynny: czynny.kontrolka.checked,
      odstepSekund: liczba(odstep.kontrolka),
      przyOdejsciu: przyOdejsciu.kontrolka.checked,
      przyZamknieciu: przyZamknieciu.kontrolka.checked,
      przyPrzelaczeniu: przyPrzelaczeniu.kontrolka.checked,
      ileKopiiZachowac: liczba(ileKopii.kontrolka),
      poIluGodzinachWygasa: liczba(ileGodzin.kontrolka),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Nastaw autozapisu nie zapisano: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    nastawy = wynik.wynik.settings;
    zdanieNastaw.textContent = kopieOpiszNastawy(nastawy);
    powiedz(`Nastawy autozapisu zapisane. ${kopieOpiszNastawy(nastawy)}`, true);
    przerysujWskaznik();
  }

  async function zapiszSamoczynnie(powod: StudioBackupReason): Promise<boolean> {
    const kod = dokument();
    if (kod === '') return false;
    zapisWToku = true;
    przerysujWskaznik();
    const wynik = await zrodlo.kontrolaAutozapisWykonaj(
      kod,
      kontekst.tresc(),
      kontekst.postac(),
      powod,
    );
    zapisWToku = false;
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Zapisu samoczynnego nie wykonano: ${wynik.blad?.message ?? ''}`, false);
      przerysujWskaznik();
      return false;
    }
    const tresc = wynik.wynik;
    if (!tresc.saved) {
      // Odpowiedź udana o zapisie nieudanym — nastawy niosą powód, wskaźnik pokaże „nieudany".
      nastawy = {
        ...(nastawy ?? { enabled: czynny.kontrolka.checked }),
        lastSaveFailed: true,
        ...(tresc.failureReason === undefined ? {} : { lastFailureReason: tresc.failureReason }),
      };
      powiedz(
        `ZAPIS SAMOCZYNNY SIĘ NIE UDAŁ: ${tresc.failureReason ?? 'rdzeń nie podał powodu'}. ` +
          'Treść leży w kopii zapasowej — nie zamykaj okna, dopóki zapis nie przejdzie.',
        false,
      );
      przerysujWskaznik();
      await odswiezKopie();
      return false;
    }
    nastawy = {
      ...(nastawy ?? { enabled: czynny.kontrolka.checked }),
      lastSaveFailed: false,
      ...(tresc.savedAt === undefined ? {} : { lastSaveAt: tresc.savedAt }),
    };
    powiedz(
      'Zapis samoczynny wykonany w OSOBNYM szeregu autozapisu' +
        (tresc.version === undefined ? '' : `, wersja ${tresc.version.id}`) +
        (tresc.backup === undefined ? '.' : `; kopia zapasowa ${tresc.backup.id} założona przed zapisem.`),
      true,
    );
    przerysujWskaznik();
    await Promise.all([odswiezKopie(), odswiezSzeregi()]);
    return true;
  }

  async function kopiaPrzedCzynnoscia(): Promise<string | null> {
    const kod = kontekst.idDokumentu();
    if (kod === '') return null;
    const wynik = await zrodlo.kontrolaKopiaZaloz(
      kod,
      StudioBackupReason.BeforeIrreversible,
      kontekst.tresc(),
      kontekst.postac(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(
        `Kopii przed czynnością nieodwracalną NIE założono: ${wynik.blad?.message ?? ''}. ` +
          'Czynność nieodwracalna bez kopii jest ryzykiem — rozważ zapis przed nią.',
        false,
      );
      return null;
    }
    await odswiezKopie();
    return wynik.wynik.backup.id;
  }

  async function zalozKopieRecznie(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaKopiaZaloz(
      kod,
      StudioBackupReason.Manual,
      kontekst.tresc(),
      kontekst.postac(),
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Kopii nie założono: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const kopia = wynik.wynik.backup;
    powiedz(
      kopia.succeeded
        ? `Kopia ${kopia.id} założona — ${kopieOpiszKopie(kopia)}`
        : `KOPIA SIĘ NIE UDAŁA: ${kopia.failureReason ?? 'rdzeń nie podał powodu'}`,
      kopia.succeeded,
    );
    await odswiezKopie();
  }

  async function odswiezKopie(): Promise<void> {
    const kod = kontekst.idDokumentu();
    const wynik = await zrodlo.kontrolaKopieWykaz(
      kod,
      kontekst.idOkna(),
      tylkoNiezapisane.kontrolka.checked,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      podsumowanieKopii.textContent = `Kopii nie udało się odczytać: ${wynik.blad?.message ?? ''}`;
      wykazKopii.replaceChildren();
      return;
    }
    kopie = wynik.wynik.backups;
    podsumowanieKopii.textContent =
      kopie.length === 0
        ? 'Kopii zapasowych nie ma. Pusty wykaz znaczy tu „nie założono jeszcze żadnej", ' +
          'a nie „kopie są wyłączone".'
        : `Kopii: ${kopie.length}; niosących zmiany niezapisane: ${wynik.wynik.unsavedCount}.`;
    wykazKopii.replaceChildren(...kopie.map(wierszKopii));

    const zgloszeniePoZamknieciu = kopieZgloszeniePoZamknieciu(kopie);
    if (zgloszeniePoZamknieciu === null) {
      zgloszenie.textContent =
        'Nie ma kopii niosącej zmiany niezapisane — nic do przywrócenia po nagłym zamknięciu.';
      zgloszenie.dataset['zgloszenie'] = 'nie';
      return;
    }
    zgloszenie.textContent = zgloszeniePoZamknieciu.zdanie;
    zgloszenie.dataset['zgloszenie'] = 'tak';
  }

  async function przywroc(kopia: StudioDocumentBackup): Promise<void> {
    const wynik = await zrodlo.kontrolaKopiaPrzywroc(
      kopia.id,
      doNowego.kontrolka.checked,
      kontekst.idOkna(),
      '',
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Przywrócenie kopii odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const tresc = wynik.wynik;
    if (doNowego.kontrolka.checked) {
      powiedz(
        `Kopia przywrócona DO NOWEGO DOKUMENTU ${tresc.document.id} — dokument stojący w oknie ` +
          'został nietknięty.',
        true,
      );
    } else {
      kontekst.naSkutek(tresc.document.content, tresc.form);
      powiedz(
        `Kopia przywrócona NA MIEJSCE w dokumencie ${tresc.document.id}. Stan sprzed przywrócenia ` +
          'stoi w kopiach, więc samo przywrócenie da się cofnąć.',
        true,
      );
    }
    await odswiezKopie();
  }

  async function odswiezSzeregi(): Promise<void> {
    const kod = kontekst.idDokumentu();
    if (kod === '') return;
    const wybrany = szereg.kontrolka.value;
    const wynik = await zrodlo.kontrolaSzeregiWersji(
      kod,
      wybrany === '' ? null : (wybrany as StudioVersionSeries),
      tylkoKluczowe.kontrolka.checked,
      100,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      podsumowanieSzeregow.textContent =
        `Szeregów wersji nie udało się odczytać: ${wynik.blad?.message ?? ''}`;
      wykazWersji.replaceChildren();
      return;
    }
    const tresc = wynik.wynik;
    podsumowanieSzeregow.textContent = kopieOpiszSzeregi(
      tresc.versions,
      tresc.operatorCount,
      tresc.autosaveCount,
      tresc.initialVersionId,
    );
    wykazWersji.replaceChildren(
      ...tresc.versions.map((wersja) =>
        wierszWersji(wersja, wybrany === '' ? undefined : (wybrany as StudioVersionSeries),
          tresc.initialVersionId),
      ),
    );
  }

  async function wrocDoZalozycielskiej(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaPowrotDoZalozycielskiej(kod, true);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Powrót do wersji założycielskiej odmówiony: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const tresc = wynik.wynik;
    kontekst.naSkutek(tresc.document.content, tresc.form);
    powiedz(
      `Dokument wrócił do wersji założycielskiej ${tresc.version.id}. Wersje nowsze ZOSTAŁY ` +
        'w historii, więc ten powrót sam da się cofnąć.',
      true,
    );
    await odswiezSzeregi();
  }

  /* ── Wiersze wykazów ────────────────────────────────────────────────────── */

  function wierszKopii(kopia: StudioDocumentBackup): HTMLElement {
    const opis = document.createElement('p');
    opis.className = 'dn-pole-opis';
    opis.textContent = kopieOpiszKopie(kopia);

    const przywrocenie = przyciskPanelu('Przywróć tę kopię', 'kopia-przywroc');
    przywrocenie.addEventListener('click', () => {
      void przywroc(kopia);
    });

    const pozycja = document.createElement('li');
    pozycja.dataset['kopia'] = kopia.id;
    pozycja.dataset['niezapisane'] = kopia.unsavedChanges === true ? 'tak' : 'nie';
    pozycja.dataset['udana'] = kopia.succeeded ? 'tak' : 'nie';
    pozycja.append(opis, przywrocenie);
    return pozycja;
  }

  function wierszWersji(
    wersja: StudioVersion,
    wybranySzereg: StudioVersionSeries | undefined,
    idZalozycielskiej: string | undefined,
  ): HTMLElement {
    const nazwa = document.createElement('p');
    nazwa.className = 'ms-kontrola__glowa';
    nazwa.textContent =
      (wersja.label === undefined || wersja.label === '' ? wersja.id : wersja.label) +
      (wersja.id === idZalozycielskiej ? ' — wersja założycielska' : '');

    const opis = document.createElement('p');
    opis.className = 'dn-pole-opis';
    opis.textContent =
      `${kopieNazwaSzeregu(wybranySzereg)} · ` +
      `${new Date(wersja.createdAt).toLocaleString('pl-PL')} · ` +
      (wersja.milestone === true ? 'kluczowa' : 'zwyczajna') +
      (wersja.summary === undefined || wersja.summary === '' ? '' : ` · ${wersja.summary}`);

    const pozycja = document.createElement('li');
    pozycja.dataset['wersja'] = wersja.id;
    pozycja.dataset['kluczowa'] = wersja.milestone === true ? 'tak' : 'nie';
    pozycja.append(nazwa, opis);
    return pozycja;
  }

  async function odswiez(): Promise<void> {
    await Promise.all([odswiezNastawy(), odswiezKopie(), odswiezSzeregi()]);
  }

  przerysujWskaznik();

  return {
    element,
    odswiez,
    zapiszSamoczynnie,
    kopiaPrzedCzynnoscia,
    stanZapisu: () => kopieStanZapisu(nastawy, kontekst.zmianyNiezapisane(), zapisWToku),
    zdanieStanuZapisu: () =>
      kopieOpiszStanZapisu(
        kopieStanZapisu(nastawy, kontekst.zmianyNiezapisane(), zapisWToku),
        nastawy,
      ),

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
      if (otwarty) void odswiez();
    },

    widoczny: () => otwarty,
  };
}

/**
 * Przycisk panelu wraz z kodem czynności do sprawdzianu, po którym testy
 * i obsługa zdarzeń rozpoznają naciśnięty przycisk.
 */
function przyciskPanelu(nazwa: string, kod: string): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  przycisk.textContent = nazwa;
  przycisk.dataset['czynnosc'] = kod;
  return przycisk;
}

/**
 * Pas przycisków — jeden rząd czynności ułożonych poziomo w jednej części
 * panelu, bez podziału na kolumny.
 */
function pasPrzyciskow(przyciski: readonly HTMLElement[]): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'ms-kontrola__pas';
  pas.append(...przyciski);
  return pas;
}

/**
 * Część panelu wraz z jej tytułem widocznym nad zebranymi w niej kontrolkami,
 * jako osobna sekcja nakładki.
 */
function czescPanelu(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-kontrola__tytul';
  naglowek.textContent = tytul;

  const sekcja = document.createElement('section');
  sekcja.className = 'ms-kontrola__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}
