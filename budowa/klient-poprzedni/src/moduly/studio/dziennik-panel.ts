import './dziennik-kontrola.css';

import {
  StudioActionState,
  StudioAuthor,
  StudioPasteMode,
  type StudioDocumentAction,
  type StudioModelChangeSummary,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import { poleLogiczne, poleTekstowe, poleWyboru } from '../../modele/kontrolki-formularza';
import {
  DZIENNIK_ROZNICA_DRZEW,
  dziennikCzyCofnieciePowstrzymane,
  dziennikOpiszBilans,
  dziennikOpiszCofniecie,
  dziennikOpiszCofniecieModelu,
  dziennikOpiszRoznicePostaci,
  dziennikOpiszRoznicePostaciWpis,
  dziennikOpiszWpis,
  dziennikOpiszZaleznosci,
  dziennikOpiszZmianyModelu,
  dziennikPoKolejnosci,
} from './dziennik-czynnosci';
import type { ZrodloKontroliStudio } from './zrodlo-kontroli-studio';

/**
 * Panel kontroli pracy nad dokumentem — cofanie, zmiany modelu, różnica postaci
 * i schowek dokumentu.
 *
 * ── Co panel niesie i czego NIE dubluje ─────────────────────────────────────
 * Historia wersji ma swoje okno (Session Repository), a zmiany śledzone swoje
 * miejsce w treści. Ten panel dokłada to, czego ani jedno, ani drugie nie ma:
 *
 *   — **dziennik czynności** — cofnięcie pojedyncze i NIE PO KOLEI, wraz
 *     z ponowieniem; zależności wypisane przy wpisie, żeby odmowa nie była
 *     zaskoczeniem;
 *   — **zmiany modelu** — licznik, skakanie i cofnięcie wszystkiego albo
 *     odhaczonych, z zachowaniem pracy Operatora;
 *   — **różnica POSTACI dwóch wersji** wraz z przeniesieniem pojedynczego
 *     fragmentu do stanu bieżącego;
 *   — **schowek dokumentu** — odłożenie i wklejenie fragmentu drogą rdzenia,
 *     która sprawdza blokady i odkłada wpis dziennika, więc wklejenie da się
 *     cofnąć pojedynczo.
 *
 * ── Dlaczego panel woła rdzeń sam ───────────────────────────────────────────
 * Bo każda z tych czynności oddaje BILANS, a bilans jest treścią dla Operatora,
 * nie wartością pośrednią: co przeszło, co stanęło i przez którą blokadę.
 * Przepuszczanie go przez okno nadrzędne kosztowałoby jedno przełożenie na
 * każdej z dziesięciu dróg, a bilans musi trafić na widok w całości. Skutek
 * czynności — nową treść i postać — panel oddaje oknu wywołaniem zwrotnym,
 * bo powierzchnia dokumentu jest po jego stronie.
 *
 * Panel nie zna treści dokumentu ani zaznaczenia: bierze je z kontekstu, który
 * podaje okno. Dzięki temu ten sam panel obsługuje dokument w zakładce i drugi
 * w podziale powierzchni.
 */

/** Czym panel pyta okno o stan dokumentu i co mu oddaje. */
export interface KontekstDziennika {
  /** Dokument czynny; puste znaczy „nie ma na czym pracować". */
  idDokumentu(): string;
  /** Zaznaczenie w treści; `null` znaczy sam kursor albo brak zaznaczenia. */
  zaznaczenie(): { poczatek: number; koniec: number } | null;
  /** Miejsce kursora w znakach treści. */
  kursor(): number;
  /** Przyjmuje treść i postać oddaną przez rdzeń po czynności. */
  naSkutek(tresc: string | undefined, postac: unknown): void;
  /** Przewija treść do wskazanego miejsca — po skoku po zmianach modelu. */
  naMiejsce(poczatek: number, koniec: number): void;
}

/** Panel wraz z czynnościami wołanymi z zewnątrz. */
export interface DziennikPanel {
  element: HTMLElement;
  /** Odczytuje dziennik i zestawienie zmian modelu. */
  odswiez(): Promise<void>;
  /** Zestawienie zmian modelu ostatnio odczytane; `null` przed pierwszym odczytem. */
  zmianyModelu(): StudioModelChangeSummary | null;
  /** Skacze do następnej albo poprzedniej zmiany modelu — dla przełącznika okna. */
  skocz(wPrzod: boolean): Promise<StudioTrackedChange | null>;
  /** Odkłada zaznaczenie do schowka dokumentu; `wytnij` wycina je z treści. */
  odlozZaznaczenie(wytnij: boolean, samaPostac: boolean): Promise<void>;
  /** Wkleja wpis schowka w miejsce kursora wedle wybranego trybu. */
  wklej(idWpisu: string, tryb: StudioPasteMode): Promise<void>;
  przestawWidocznosc(): void;
  widoczny(): boolean;
}

export function utworzDziennikPanel(
  zrodlo: ZrodloKontroliStudio,
  kontekst: KontekstDziennika,
): DziennikPanel {
  let otwarty = false;
  let wpisy: readonly StudioDocumentAction[] = [];
  let zestawienie: StudioModelChangeSummary | null = null;
  const odhaczone = new Set<string>();
  let ostatnieMiejsce = 0;

  /* ── Odpowiedź rdzenia; cisza jest zakazana ──────────────────────────────── */

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-pole-opis ms-kontrola__odpowiedz';
  odpowiedz.dataset['czynnosc'] = 'odpowiedz';

  function powiedz(tresc: string, udana: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.dataset['udana'] = udana ? 'tak' : 'nie';
  }

  /** Dokument czynny albo odpowiedź nazywająca brak; `''` znaczy „nie ma na czym". */
  function dokument(): string {
    const kod = kontekst.idDokumentu();
    if (kod === '') {
      powiedz(
        'Nie ma dokumentu czynnego. Dziennik czynności, zmiany modelu i schowek dokumentu ' +
          'dotyczą jednego dokumentu — otwórz go albo załóż nowy.',
        false,
      );
    }
    return kod;
  }

  /* ── Dziennik czynności ──────────────────────────────────────────────────── */

  const tylkoStojace = poleLogiczne({
    etykieta: 'Tylko czynności stojące w dokumencie',
    opis:
      'Zdjęcie zawężenia pokazuje także czynności cofnięte — te da się PONOWIĆ. Cofnięta ' +
      'czynność nie znika z dziennika, bo cofnięcie samo ma być odwracalne.',
  });
  tylkoStojace.kontrolka.checked = true;
  tylkoStojace.kontrolka.addEventListener('change', () => {
    void odswiezDziennik();
  });

  const tylkoModel = poleLogiczne({
    etykieta: 'Tylko czynności modelu',
    opis: 'Zawężenie idzie polem autora czynności w rdzeniu, nie odsiewem w oknie.',
  });
  tylkoModel.kontrolka.addEventListener('change', () => {
    void odswiezDziennik();
  });

  const podsumowanieDziennika = document.createElement('p');
  podsumowanieDziennika.className = 'dn-pole-opis';

  const zdanieRoznicyDrzew = document.createElement('p');
  zdanieRoznicyDrzew.className = 'dn-pole-opis ms-kontrola__zasada';
  zdanieRoznicyDrzew.textContent = DZIENNIK_ROZNICA_DRZEW;

  const wykazCzynnosci = document.createElement('ul');
  wykazCzynnosci.className = 'ms-kontrola__wykaz';

  /* ── Zmiany modelu ──────────────────────────────────────────────────────── */

  const licznikZmian = document.createElement('p');
  licznikZmian.className = 'dn-pole-opis ms-kontrola__licznik';
  licznikZmian.dataset['czynnosc'] = 'licznik-zmian-modelu';
  licznikZmian.textContent = 'Zmiany modelu nie zostały jeszcze odczytane z rdzenia.';

  const tylkoPostac = poleLogiczne({
    etykieta: 'Tylko zmiany postaci modelu',
    opis:
      'Zmiana kroju albo wcięcia bez zmiany liter jest zmianą modelu tak samo jak dopisany ' +
      'akapit — to zawężenie pokazuje wyłącznie ją.',
  });
  tylkoPostac.kontrolka.addEventListener('change', () => {
    void odswiezZmianyModelu();
  });

  const wstecz = przyciskPanelu('Poprzednia zmiana modelu', 'skok-wstecz');
  wstecz.addEventListener('click', () => {
    void skocz(false);
  });
  const wprzod = przyciskPanelu('Następna zmiana modelu', 'skok-wprzod');
  wprzod.addEventListener('click', () => {
    void skocz(true);
  });

  const cofnijWszystkie = przyciskPanelu('Cofnij WSZYSTKO, co zrobił model', 'cofnij-model-wszystko');
  cofnijWszystkie.title =
    'Cofa zmiany modelu z ZACHOWANIEM zmian Operatora naniesionych w tym czasie. Nie jest to ' +
    'przywrócenie wersji sprzed — to skasowałoby też pracę Operatora. Kopia zapasowa zakładana ' +
    'jest PRZED cofnięciem.';
  cofnijWszystkie.addEventListener('click', () => {
    void cofnijModel(true);
  });

  const cofnijOdhaczone = przyciskPanelu('Cofnij odhaczone zmiany modelu', 'cofnij-model-wybrane');
  cofnijOdhaczone.title =
    'Cofa wyłącznie zmiany odhaczone w wykazie poniżej — jednym poleceniem, wykazem kodów ' +
    'przyjmowanym przez rdzeń.';
  cofnijOdhaczone.addEventListener('click', () => {
    void cofnijModel(false);
  });

  const wykazZmianModelu = document.createElement('ul');
  wykazZmianModelu.className = 'ms-kontrola__wykaz';

  /* ── Różnica postaci i przeniesienie fragmentu ──────────────────────────── */

  const wersjaOdniesienia = poleTekstowe({
    etykieta: 'Wersja odniesienia różnicy postaci',
    podpowiedz: 'puste = wersja założycielska dokumentu',
    opis: 'Pole puste czyta rdzeń jako wersję założycielską, a nie jako brak wskazania.',
  });
  const wersjaPorownywana = poleTekstowe({
    etykieta: 'Wersja porównywana',
    podpowiedz: 'puste = stan bieżący dokumentu',
  });
  const obszarPostaci = poleWyboru(
    {
      etykieta: 'Obszar postaci',
      opis: 'Zawężenie do jednego obszaru; „wszystkie" oddaje różnicę całej postaci.',
    },
    [
      { wartosc: '', etykieta: 'Wszystkie obszary postaci' },
      { wartosc: 'characterStyle', etykieta: 'Styl znaku' },
      { wartosc: 'paragraphStyle', etykieta: 'Styl akapitu' },
      { wartosc: 'namedStyle', etykieta: 'Style nazwane' },
      { wartosc: 'section', etykieta: 'Sekcje i nastawy strony' },
      { wartosc: 'table', etykieta: 'Tabele' },
      { wartosc: 'object', etykieta: 'Obiekty osadzone' },
      { wartosc: 'apparatus', etykieta: 'Aparat dokumentu' },
    ],
  );

  const porownajPostac = przyciskPanelu('Porównaj postać dwóch wersji', 'roznica-postaci');
  porownajPostac.title =
    'studio.diff.form.compare — zmiana kroju czy wcięcia jest widoczna jako zmiana, a nie milczy, ' +
    'bo litery zostały te same.';
  porownajPostac.addEventListener('click', () => {
    void porownaj();
  });

  const podsumowaniePostaci = document.createElement('p');
  podsumowaniePostaci.className = 'dn-pole-opis';

  const wykazPostaci = document.createElement('ul');
  wykazPostaci.className = 'ms-kontrola__wykaz';

  const numerFragmentu = poleTekstowe({
    etykieta: 'Numer fragmentu różnicy do przeniesienia',
    podpowiedz: 'numer fragmentu z widoku różnicy',
    opis:
      'Przeniesienie bierze fragment ze wskazanej wersji i wnosi go do stanu BIEŻĄCEGO — to jest ' +
      'praktyczny sens widoku różnicy, nie samo patrzenie.',
  });
  const zPostacia = poleLogiczne({ etykieta: 'Przenieś także postać fragmentu' });
  zPostacia.kontrolka.checked = true;

  const przenies = przyciskPanelu('Przenieś fragment do stanu bieżącego', 'przenies-fragment');
  przenies.addEventListener('click', () => {
    void przeniesFragment();
  });

  /* ── Schowek dokumentu ──────────────────────────────────────────────────── */

  const idWpisuSchowka = poleTekstowe({
    etykieta: 'Wpis schowka do wklejenia',
    podpowiedz: 'puste = wpis najświeższy',
    opis:
      'Droga rdzenia sprawdza blokady fragmentu PRZED dotknięciem treści i odkłada wpis ' +
      'dziennika, więc wklejenie da się cofnąć pojedynczo. Rodzina clipboard.* platformy tego ' +
      'nie robi i dlatego są to dwie różne czynności, nie dwie nazwy jednej.',
  });

  const trybWklejenia = poleWyboru(
    {
      etykieta: 'Sposób wklejenia',
      opis: 'Wybór należy do Operatora i jest jawny — żaden z trzech nie jest ukrytym domyślnym.',
    },
    [
      { wartosc: StudioPasteMode.KeepFormat, etykieta: 'Z zachowaniem postaci źródła' },
      { wartosc: StudioPasteMode.PlainText, etykieta: 'Jako czysty tekst' },
      { wartosc: StudioPasteMode.MergeFormat, etykieta: 'Z przejęciem postaci miejsca wklejenia' },
    ],
  );

  const odlozKopie = przyciskPanelu('Odłóż zaznaczenie do schowka', 'schowek-odloz');
  odlozKopie.addEventListener('click', () => {
    void odlozZaznaczenie(false, false);
  });
  const odlozWyciecie = przyciskPanelu('Wytnij zaznaczenie do schowka', 'schowek-wytnij');
  odlozWyciecie.addEventListener('click', () => {
    void odlozZaznaczenie(true, false);
  });
  const odlozPostac = przyciskPanelu('Zabierz samą postać (malarz formatów)', 'schowek-postac');
  odlozPostac.title =
    'Odkłada POSTAĆ fragmentu, nie jego treść — tak działa malarz formatów. Treść dokumentu ' +
    'zostaje nietknięta.';
  odlozPostac.addEventListener('click', () => {
    void odlozZaznaczenie(false, true);
  });
  const wklejPrzycisk = przyciskPanelu('Wklej w miejsce kursora', 'schowek-wklej');
  wklejPrzycisk.addEventListener('click', () => {
    void wklej(idWpisuSchowka.kontrolka.value.trim(), trybWklejenia.kontrolka.value as StudioPasteMode);
  });

  /* ── Złożenie panelu ────────────────────────────────────────────────────── */

  const element = document.createElement('section');
  element.className = 'ms-kontrola';
  element.dataset['panel'] = 'kontrola-pracy';
  element.hidden = true;
  element.setAttribute('aria-label', 'Kontrola pracy nad dokumentem — cofanie, zmiany modelu, schowek');
  element.append(
    czescPanelu('Dziennik czynności — cofanie pojedyncze, także nie po kolei', [
      zdanieRoznicyDrzew,
      tylkoStojace.element,
      tylkoModel.element,
      podsumowanieDziennika,
      wykazCzynnosci,
    ]),
    czescPanelu('Wszystko, co zrobił model', [
      licznikZmian,
      tylkoPostac.element,
      pasPrzyciskow([wstecz, wprzod, cofnijWszystkie, cofnijOdhaczone]),
      wykazZmianModelu,
    ]),
    czescPanelu('Różnica POSTACI dwóch wersji i przeniesienie fragmentu', [
      wersjaOdniesienia.element,
      wersjaPorownywana.element,
      obszarPostaci.element,
      pasPrzyciskow([porownajPostac]),
      podsumowaniePostaci,
      wykazPostaci,
      numerFragmentu.element,
      zPostacia.element,
      pasPrzyciskow([przenies]),
    ]),
    czescPanelu('Schowek dokumentu — drogą rdzenia, ze sprawdzeniem blokad', [
      idWpisuSchowka.element,
      trybWklejenia.element,
      pasPrzyciskow([odlozKopie, odlozWyciecie, odlozPostac, wklejPrzycisk]),
    ]),
    odpowiedz,
  );

  /* ── Czynności ──────────────────────────────────────────────────────────── */

  async function odswiezDziennik(): Promise<void> {
    const kod = kontekst.idDokumentu();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaDziennik(kod, {
      ...(tylkoStojace.kontrolka.checked ? { stan: StudioActionState.Active } : {}),
      ...(tylkoModel.kontrolka.checked ? { autor: StudioAuthor.Model } : {}),
      granica: 200,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      podsumowanieDziennika.textContent =
        `Dziennika nie udało się odczytać: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`;
      wykazCzynnosci.replaceChildren();
      return;
    }
    wpisy = dziennikPoKolejnosci(wynik.wynik.actions);
    podsumowanieDziennika.textContent =
      wynik.wynik.total === 0
        ? 'Dziennik jest pusty: na tym dokumencie nie wykonano jeszcze ani jednej czynności ' +
          'odkładanej w dzienniku. Pusty wykaz nie znaczy tu „nie da się cofać".'
        : `Czynności w dzienniku: ${wynik.wynik.total}; w wykazie po zawężeniu: ${wpisy.length}.`;
    wykazCzynnosci.replaceChildren(...wpisy.map(wierszCzynnosci));
  }

  async function odswiezZmianyModelu(): Promise<void> {
    const kod = kontekst.idDokumentu();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaZmianyModelu(kod, false, tylkoPostac.kontrolka.checked);
    if (!wynik.udany || wynik.wynik === undefined) {
      licznikZmian.textContent =
        `Zmian modelu nie udało się odczytać: ${wynik.blad?.message ?? 'rdzeń nie podał powodu'}`;
      wykazZmianModelu.replaceChildren();
      return;
    }
    zestawienie = wynik.wynik.summary;
    licznikZmian.textContent = dziennikOpiszZmianyModelu(zestawienie);
    const zmiany = zestawienie.changes ?? [];
    wykazZmianModelu.replaceChildren(...zmiany.map(wierszZmianyModelu));
  }

  async function skocz(wPrzod: boolean): Promise<StudioTrackedChange | null> {
    const kod = dokument();
    if (kod === '') return null;
    const wynik = await zrodlo.kontrolaSkokPoZmianachModelu(
      kod,
      ostatnieMiejsce,
      wPrzod,
      tylkoPostac.kontrolka.checked,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Skok po zmianach modelu odmówiony: ${wynik.blad?.message ?? ''}`, false);
      return null;
    }
    const zmiana = wynik.wynik.change;
    if (zmiana === undefined) {
      powiedz(
        `Nie ma dalszej zmiany modelu w tym kierunku. Zmian modelu w dokumencie: ${wynik.wynik.total}.`,
        true,
      );
      return null;
    }
    ostatnieMiejsce = wPrzod ? zmiana.rangeEnd : zmiana.rangeStart;
    kontekst.naMiejsce(zmiana.rangeStart, zmiana.rangeEnd);
    const ktora = wynik.wynik.index === undefined ? '' : ` (${wynik.wynik.index} z ${wynik.wynik.total})`;
    powiedz(`Zmiana modelu na znakach ${zmiana.rangeStart}–${zmiana.rangeEnd}${ktora}.`, true);
    return zmiana;
  }

  async function cofnijModel(wszystkie: boolean): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const kody = wszystkie ? [] : [...odhaczone];
    if (!wszystkie && kody.length === 0) {
      powiedz(
        'Nie odhaczono ani jednej zmiany. Odhacz zmiany do cofnięcia w wykazie poniżej albo ' +
          'użyj cofnięcia wszystkiego.',
        false,
      );
      return;
    }
    const wynik = await zrodlo.kontrolaCofnijZmianyModelu(kod, kody, wszystkie);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Cofnięcie zmian modelu odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const tresc = wynik.wynik;
    powiedz(
      dziennikOpiszCofniecieModelu(
        tresc.revertedCount,
        tresc.keptOperatorChanges,
        tresc.balance,
        tresc.backupId,
      ),
      tresc.revertedCount > 0,
    );
    odhaczone.clear();
    kontekst.naSkutek(tresc.document.content, tresc.form);
    await odswiez();
  }

  async function cofnij(wpis: StudioDocumentAction): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaCofnijCzynnosci(kod, [wpis.id], true);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Cofnięcie odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const ocena = dziennikOpiszCofniecie(wynik.wynik.reverted, wynik.wynik.balance);
    powiedz(ocena.zdanie, ocena.udane);
    kontekst.naSkutek(wynik.wynik.document.content, wynik.wynik.form);
    await odswiez();
  }

  async function ponow(wpis: StudioDocumentAction): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaPonowCzynnosci(kod, [wpis.id]);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Ponowienie odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    powiedz(
      `Ponowiono czynności: ${wynik.wynik.redone.length}. ${dziennikOpiszBilans(wynik.wynik.balance)}`,
      wynik.wynik.redone.length > 0,
    );
    kontekst.naSkutek(wynik.wynik.document.content, wynik.wynik.form);
    await odswiez();
  }

  async function porownaj(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const wynik = await zrodlo.kontrolaRoznicaPostaci(
      kod,
      wersjaOdniesienia.kontrolka.value.trim(),
      wersjaPorownywana.kontrolka.value.trim(),
      obszarPostaci.kontrolka.value,
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      podsumowaniePostaci.textContent =
        `Różnicy postaci nie udało się policzyć: ${wynik.blad?.message ?? ''}`;
      wykazPostaci.replaceChildren();
      return;
    }
    const tresc = wynik.wynik;
    podsumowaniePostaci.textContent = dziennikOpiszRoznicePostaci(tresc.entries, tresc);
    wykazPostaci.replaceChildren(
      ...tresc.entries.map((wpis) => {
        const pozycja = document.createElement('li');
        pozycja.dataset['obszar'] = wpis.area;
        pozycja.dataset['rodzaj'] = wpis.kind;
        pozycja.textContent = dziennikOpiszRoznicePostaciWpis(wpis);
        return pozycja;
      }),
    );
  }

  async function przeniesFragment(): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const zrodloWersji = wersjaOdniesienia.kontrolka.value.trim();
    if (zrodloWersji === '') {
      powiedz(
        'Przeniesienie fragmentu wymaga wskazania wersji, Z KTÓREJ fragment ma przyjść: wpisz ją ' +
          'w polu wersji odniesienia. Bez niej rdzeń nie wie, skąd brać brzmienie.',
        false,
      );
      return;
    }
    const numer = Number(numerFragmentu.kontrolka.value.trim());
    const zaznaczenie = kontekst.zaznaczenie();
    const wskazanie = Number.isFinite(numer) && numerFragmentu.kontrolka.value.trim() !== ''
      ? { numerFragmentu: numer }
      : zaznaczenie === null
        ? {}
        : { poczatek: zaznaczenie.poczatek, koniec: zaznaczenie.koniec };
    if (Object.keys(wskazanie).length === 0) {
      powiedz(
        'Wskaż fragment: podaj jego numer z widoku różnicy albo zaznacz zakres w treści. ' +
          'Bez wskazania rdzeń nie wie, który fragment przenieść.',
        false,
      );
      return;
    }
    const wynik = await zrodlo.kontrolaPrzeniesFragment(
      kod,
      zrodloWersji,
      wskazanie,
      zPostacia.kontrolka.checked,
      {},
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Przeniesienie fragmentu odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    powiedz(
      `Fragment przeniesiony z wersji ${zrodloWersji}. ${dziennikOpiszBilans(wynik.wynik.balance)}`,
      wynik.wynik.balance.applied > 0,
    );
    kontekst.naSkutek(wynik.wynik.document.content, wynik.wynik.form);
    await odswiez();
  }

  async function odlozZaznaczenie(wytnij: boolean, samaPostac: boolean): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const zaznaczenie = kontekst.zaznaczenie();
    if (zaznaczenie === null || zaznaczenie.poczatek === zaznaczenie.koniec) {
      powiedz(
        'Nic nie jest zaznaczone. Schowek dokumentu odkłada FRAGMENT, więc bez zaznaczenia nie ' +
          'ma czego odłożyć — zaznacz fragment w treści.',
        false,
      );
      return;
    }
    const wynik = await zrodlo.kontrolaSchowekOdlozFragment(
      kod,
      zaznaczenie,
      wytnij,
      samaPostac,
      {},
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Odłożenie do schowka odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    const tresc = wynik.wynik;
    idWpisuSchowka.kontrolka.value = tresc.clipboardEntryId;
    powiedz(
      `Odłożono do schowka wpis ${tresc.clipboardEntryId} — ${tresc.text.length} znaków` +
        (samaPostac ? ' (sama postać, treść dokumentu nietknięta)' : '') +
        (wytnij ? ' i wycięto fragment z treści' : '') +
        (tresc.actionId === undefined
          ? '.'
          : `. Czynność stoi w dzienniku jako ${tresc.actionId}, więc da się ją cofnąć pojedynczo.`),
      true,
    );
    if (wytnij) kontekst.naSkutek(undefined, tresc.form);
    await odswiez();
  }

  async function wklej(idWpisu: string, tryb: StudioPasteMode): Promise<void> {
    const kod = dokument();
    if (kod === '') return;
    const zaznaczenie = kontekst.zaznaczenie();
    const wynik = await zrodlo.kontrolaSchowekWklejFragment(
      kod,
      kontekst.kursor(),
      { idWpisu },
      tryb,
      // Zaznaczenie niepuste znaczy „wklej ZAMIAST tego": tak działa wklejenie
      // w pakiecie biurowym i inne zachowanie byłoby tu zaskoczeniem.
      zaznaczenie === null || zaznaczenie.poczatek === zaznaczenie.koniec
        ? {}
        : { poczatek: zaznaczenie.poczatek, koniec: zaznaczenie.koniec },
      {},
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(`Wklejenie odmówione: ${wynik.blad?.message ?? ''}`, false);
      return;
    }
    powiedz(
      `Wklejono w miejscu ${kontekst.kursor()}. ${dziennikOpiszBilans(wynik.wynik.balance)}`,
      wynik.wynik.balance.applied > 0,
    );
    kontekst.naSkutek(undefined, wynik.wynik.form);
    await odswiez();
  }

  /* ── Wiersze wykazów ────────────────────────────────────────────────────── */

  function wierszCzynnosci(wpis: StudioDocumentAction): HTMLElement {
    const glowa = document.createElement('p');
    glowa.className = 'ms-kontrola__glowa';
    glowa.textContent = wpis.description === '' ? dziennikOpiszWpis(wpis) : wpis.description;

    const opis = document.createElement('p');
    opis.className = 'dn-pole-opis';
    opis.textContent = dziennikOpiszWpis(wpis);

    const zaleznosci = document.createElement('p');
    zaleznosci.className = 'dn-pole-opis';
    zaleznosci.textContent = dziennikOpiszZaleznosci(wpis);

    const pas = document.createElement('div');
    pas.className = 'ms-kontrola__pas';

    if (wpis.state === StudioActionState.Reverted) {
      const ponowienie = przyciskPanelu('Ponów tę czynność', 'ponow');
      ponowienie.addEventListener('click', () => {
        void ponow(wpis);
      });
      pas.append(ponowienie);
    } else {
      const cofniecie = przyciskPanelu('Cofnij tę czynność', 'cofnij');
      cofniecie.title = DZIENNIK_ROZNICA_DRZEW;
      // Przycisk zostaje CZYNNY nawet przy zależności: prawdę rozstrzyga rdzeń,
      // a wykaz zależności mógł się zmienić po ostatnim odczycie. Operator
      // dowiaduje się o przeszkodzie zdaniem obok, a odmowę nazywa rdzeń.
      cofniecie.dataset['zaleznosc'] = dziennikCzyCofnieciePowstrzymane(wpis) ? 'tak' : 'nie';
      cofniecie.addEventListener('click', () => {
        void cofnij(wpis);
      });
      pas.append(cofniecie);
    }

    const pozycja = document.createElement('li');
    pozycja.dataset['czynnoscDziennika'] = wpis.id;
    pozycja.dataset['stan'] = wpis.state;
    pozycja.dataset['autor'] = wpis.author;
    pozycja.append(glowa, opis, zaleznosci, pas);
    return pozycja;
  }

  function wierszZmianyModelu(zmiana: StudioTrackedChange): HTMLElement {
    const odhacz = document.createElement('input');
    odhacz.type = 'checkbox';
    odhacz.className = 'dn-przelacznik';
    odhacz.dataset['zmiana'] = zmiana.id;
    odhacz.checked = odhaczone.has(zmiana.id);
    odhacz.setAttribute('aria-label', `Odhacz zmianę modelu ${zmiana.id} do cofnięcia`);
    odhacz.addEventListener('change', () => {
      if (odhacz.checked) odhaczone.add(zmiana.id);
      else odhaczone.delete(zmiana.id);
    });

    const opis = document.createElement('span');
    opis.textContent =
      `${zmiana.kind} · znaki ${zmiana.rangeStart}–${zmiana.rangeEnd} · ` +
      `przed ${(zmiana.before ?? '').length}, po ${(zmiana.after ?? '').length} znaków · ` +
      `stan decyzji: ${zmiana.decision}`;

    const doMiejsca = przyciskPanelu('Pokaż miejsce', 'pokaz-miejsce');
    doMiejsca.addEventListener('click', () => {
      ostatnieMiejsce = zmiana.rangeStart;
      kontekst.naMiejsce(zmiana.rangeStart, zmiana.rangeEnd);
    });

    const pozycja = document.createElement('li');
    pozycja.dataset['zmianaModelu'] = zmiana.id;
    pozycja.append(odhacz, ' ', opis, ' ', doMiejsca);
    return pozycja;
  }

  async function odswiez(): Promise<void> {
    await Promise.all([odswiezDziennik(), odswiezZmianyModelu()]);
  }

  return {
    element,
    odswiez,
    zmianyModelu: () => zestawienie,
    skocz,
    odlozZaznaczenie,
    wklej,

    przestawWidocznosc() {
      otwarty = !otwarty;
      element.hidden = !otwarty;
      if (otwarty) void odswiez();
    },

    widoczny: () => otwarty,
  };
}

/** Przycisk panelu wraz z kodem czynności do sprawdzianu. */
function przyciskPanelu(nazwa: string, kod: string): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  przycisk.textContent = nazwa;
  przycisk.dataset['czynnosc'] = kod;
  return przycisk;
}

/** Pas przycisków — jeden rząd czynności. */
function pasPrzyciskow(przyciski: readonly HTMLElement[]): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'ms-kontrola__pas';
  pas.append(...przyciski);
  return pas;
}

/** Część panelu wraz z jej tytułem. */
function czescPanelu(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-kontrola__tytul';
  naglowek.textContent = tytul;

  const sekcja = document.createElement('section');
  sekcja.className = 'ms-kontrola__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}
