import {
  StudioExportFormat,
  StudioImportFormat,
  StudioTemplateFieldKind,
  type StudioTemplateDetail,
  type StudioTemplateFieldSpec,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLiczbowe,
  poleTekstowe,
  poleTresci,
  przycisk,
  utworzWierszOdpowiedzi,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { Wynik } from '../../protokol/kanal';
import type { StanStudio } from './stan-studio';
import type { SzablonZrodlo } from './szablon-zrodlo';
import {
  wstawieniaOpiszBilans,
  wstawieniaOpiszWniesienie,
  wstawieniaOpiszWydanie,
} from './zrodlo-wstawien-studio';

/**
 * Warsztat szablonów pism — nakładka na żądanie.
 *
 * ── Czego brakowało, choć galeria stała ─────────────────────────────────────
 * Galeria (`galeria-szablonow.ts`) pokazywała wykaz FABRYCZNY i zakładała z niego
 * dokument. Nie było czym szablonu założyć, zmienić, wnieść z pliku Operatora ani
 * oddać do pliku — czyli warsztatu nie było wcale. Ten panel go zamyka siedmioma
 * komendami rodziny `studio.template.*`.
 *
 * ── Szablon niesie wszystkie rzeczy naraz ───────────────────────────────────
 * Wymaganie Właściciela wprost: szablon zapisany z bieżącego dokumentu bierze
 * arkusz stylów, nastawy strony, nagłówek, stopkę, logo, tabele I BLOKADY
 * WZORCOWE. Blokady idą polem `includeLocks`, którego brak znaczy „tak" —
 * fragmenty wzorcowe pisma mają zostać wzorcowe także w dokumentach z szablonu.
 *
 * ── Szablonu fabrycznego się nie usuwa ──────────────────────────────────────
 * Rdzeń odmawia i panel czyta pole `deleted`: odpowiedź „nie usunąłem" nazywa
 * powód i pozycja zostaje w wykazie. Zdjęcie jej z widoku i pozwolenie, żeby
 * wróciła przy następnym odczycie, byłoby udawaniem skutku.
 */
export interface SzablonPanel {
  element: HTMLElement;
  przestawWidocznosc(): void;
  odswiez(): void;
  /** Wskazuje szablon z zewnątrz — woła to galeria po naciśnięciu miniatury. */
  wskaz(idSzablonu: string): void;
}

/** Cztery rodzaje pola szablonu z kontraktu. */
const RODZAJE_POL: readonly (readonly [string, string])[] = [
  [StudioTemplateFieldKind.Text, 'tekst'],
  [StudioTemplateFieldKind.Date, 'data'],
  [StudioTemplateFieldKind.Number, 'liczba'],
  [StudioTemplateFieldKind.Choice, 'wybór z wykazu'],
];

/** Formaty pliku szablonu Operatora — wnoszone i oddawane. */
const FORMATY_WNOSZONE: readonly (readonly [string, string])[] = [
  ['', 'rozpoznanie po zawartości'],
  [StudioImportFormat.Dotx, 'szablon Word (dotx)'],
  [StudioImportFormat.Ott, 'szablon OpenDocument (ott)'],
  [StudioImportFormat.Docx, 'dokument Word (docx)'],
  [StudioImportFormat.Odt, 'dokument OpenDocument (odt)'],
];

const FORMATY_ODDAWANE: readonly (readonly [string, string])[] = [
  ['', 'postać własna platformy'],
  [StudioExportFormat.Docx, 'dokument Word (docx)'],
  [StudioExportFormat.Odt, 'dokument OpenDocument (odt)'],
];

export function utworzSzablonPanel(stan: StanStudio, zrodlo: SzablonZrodlo): SzablonPanel {
  const odpowiedz = utworzWierszOdpowiedzi();

  let wskazany = '';
  let szczegoly: StudioTemplateDetail | null = null;
  let polaSzablonu: readonly StudioTemplateFieldSpec[] = [];

  /* ── Założenie szablonu z bieżącego dokumentu ──────────────────────────── */

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa szablonu',
    podpowiedz: 'np. pismo procesowe — wzór kancelarii',
  });
  const przeznaczenie = poleTresci('Przeznaczenie szablonu', 2, 'do czego ten wzór służy');
  const kategoria = poleTekstowe({
    etykieta: 'Kategoria szablonu',
    podpowiedz: 'czym Operator porządkuje galerię',
  });
  const zBlokadami = document.createElement('input');
  zBlokadami.type = 'checkbox';
  zBlokadami.className = 'dn-przelacznik';
  zBlokadami.checked = true;
  zBlokadami.setAttribute('aria-label', 'Przenieś blokady fragmentów do szablonu');
  const etykietaBlokad = document.createElement('label');
  etykietaBlokad.className = 'ms-szablon__zawezenie';
  etykietaBlokad.append(
    zBlokadami,
    document.createTextNode(
      'przenieś blokady fragmentów — fragmenty wzorcowe zostaną wzorcowe w dokumentach z szablonu',
    ),
  );

  const zapisz = przycisk(
    'Zapisz szablon z bieżącego dokumentu',
    'dn-btn dn-btn--sm dn-btn--sygnal',
  );
  zapisz.dataset['czynnosc'] = 'zapisz-szablon';
  zapisz.title =
    'studio.template.save bierze z dokumentu arkusz stylów, nastawy strony, nagłówek, stopkę, logo ' +
    'i tabele. Wskazanie szablonu w polu niżej zmienia szablon istniejący; puste zakłada nowy.';
  zapisz.addEventListener('click', () => void zapiszSzablon());

  const wskazanieSzablonu = poleTekstowe({
    etykieta: 'Szablon zmieniany',
    podpowiedz: 'puste zakłada nowy szablon',
  });
  wskazanieSzablonu.kontrolka.addEventListener('change', () => {
    wskazany = wskazanieSzablonu.kontrolka.value.trim();
    if (wskazany !== '') void odczytajPola();
  });

  /* ── Pola do wypełnienia ───────────────────────────────────────────────── */

  const nazwaPola = poleTekstowe({
    etykieta: 'Nazwa pola używana w treści szablonu',
    podpowiedz: 'np. adresat',
  });
  const etykietaPola = poleTekstowe({
    etykieta: 'Etykieta pola widoczna dla Operatora',
    podpowiedz: 'np. Adresat pisma',
  });
  const opisPola = poleTekstowe({ etykieta: 'Opis pola', podpowiedz: '' });
  const rodzajPola = wybor('Rodzaj pola', RODZAJE_POL.map((pozycja) => pozycja));
  const wartoscDomyslna = poleTekstowe({ etykieta: 'Wartość domyślna', podpowiedz: '' });
  const wyborWartosci = poleTekstowe({
    etykieta: 'Wykaz wartości do wyboru',
    podpowiedz: 'rozdzielone przecinkiem; dotyczy rodzaju „wybór z wykazu"',
  });
  const miejscePola = poleLiczbowe('Miejsce pola w treści szablonu w znakach', '');
  const wymagane = document.createElement('input');
  wymagane.type = 'checkbox';
  wymagane.className = 'dn-przelacznik';
  wymagane.setAttribute('aria-label', 'Pole jest wymagane');
  const etykietaWymagania = document.createElement('label');
  etykietaWymagania.className = 'ms-szablon__zawezenie';
  etykietaWymagania.append(wymagane, document.createTextNode('pole jest wymagane'));

  const ustawPole = przycisk('Zapisz pole szablonu', 'dn-btn dn-btn--sm dn-btn--atrament');
  ustawPole.dataset['czynnosc'] = 'ustaw-pole-szablonu';
  ustawPole.addEventListener('click', () => void ustawPoleSzablonu(false));

  const usunPole = przycisk('Usuń pole szablonu', 'dn-btn dn-btn--sm dn-btn--duch');
  usunPole.dataset['czynnosc'] = 'usun-pole-szablonu';
  usunPole.addEventListener('click', () => void ustawPoleSzablonu(true));

  const wykazPol = document.createElement('ul');
  wykazPol.className = 'ms-szablon__pola';
  wykazPol.setAttribute('aria-label', 'Pola szablonu do wypełnienia');

  /* ── Wypełnienie szablonu ──────────────────────────────────────────────── */

  const tytulDokumentu = poleTekstowe({
    etykieta: 'Tytuł dokumentu zakładanego z szablonu',
    podpowiedz: '',
  });
  const formularz = document.createElement('div');
  formularz.className = 'ms-szablon__formularz';
  const kontrolkiWartosci = new Map<string, HTMLInputElement>();

  const wypelnij = przycisk('Wypełnij szablon i oddaj dokument', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wypelnij.dataset['czynnosc'] = 'wypelnij-szablon';
  wypelnij.addEventListener('click', () => void wypelnijSzablon());

  /* ── Plik Operatora ───────────────────────────────────────────────────── */

  const sciezkaWniesienia = poleTekstowe({
    etykieta: 'Ścieżka pliku szablonu widziana przez rdzeń',
    podpowiedz: 'dotx albo ott',
    opis:
      'Klient dysku nie czyta — podaje wskazanie, a plik wciąga rdzeń. Plik Biblioteki wnosi się ' +
      'tym samym polem, wpisując jego identyfikator i wybierając rozpoznanie po zawartości.',
  });
  const plikBiblioteki = poleTekstowe({
    etykieta: 'Plik Biblioteki Library',
    podpowiedz: 'identyfikator pliku',
  });
  const formatWnoszony = wybor('Format pliku wnoszonego', FORMATY_WNOSZONE.map((p) => p));
  const wnies = przycisk('Wnieś szablon z pliku', 'dn-btn dn-btn--sm dn-btn--atrament');
  wnies.dataset['czynnosc'] = 'wnies-szablon';
  wnies.addEventListener('click', () => void wniesSzablon());

  const sciezkaOddania = poleTekstowe({
    etykieta: 'Ścieżka pliku, do którego oddać szablon',
    podpowiedz: 'puste zostawia wynik zasobem magazynu rdzenia',
  });
  const formatOddawany = wybor('Format pliku szablonu', FORMATY_ODDAWANE.map((p) => p));
  const oddaj = przycisk('Oddaj szablon do pliku', 'dn-btn dn-btn--sm dn-btn--zarys');
  oddaj.dataset['czynnosc'] = 'oddaj-szablon';
  oddaj.addEventListener('click', () => void oddajSzablon());

  const usunSzablon = przycisk('Usuń szablon własny', 'dn-btn dn-btn--sm dn-btn--duch');
  usunSzablon.dataset['czynnosc'] = 'usun-szablon';
  usunSzablon.title =
    'Szablonu FABRYCZNEGO rdzeń nie usuwa i odpowiada odmową nazywającą powód — tak samo jak ' +
    'studio.operation.delete dla operacji fabrycznych.';
  usunSzablon.addEventListener('click', () => void usunSzablonWlasny());

  const opisSzablonu = document.createElement('p');
  opisSzablonu.className = 'dn-pole-opis ms-szablon__opis';

  /* ── Czynności ─────────────────────────────────────────────────────────── */

  function idSzablonu(): string | null {
    if (wskazany === '') {
      odpowiedz.pokaz(BRAK_SZABLONU, false);
      return null;
    }
    return wskazany;
  }

  function przyjalSie(nazwaCzynnosci: string, wynik: Wynik<unknown>): boolean {
    if (wynik.udany && wynik.wynik !== undefined) return true;
    odpowiedz.pokaz(opisOdmowyBledu(nazwaCzynnosci, wynik.blad), false);
    return false;
  }

  async function zapiszSzablon(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      odpowiedz.pokaz(BRAK_DOKUMENTU, false);
      return;
    }
    if (nazwa.kontrolka.value.trim() === '') {
      odpowiedz.pokaz(BRAK_NAZWY, false);
      return;
    }
    const zmieniany = wskazanieSzablonu.kontrolka.value.trim();
    const wynik = await zrodlo.zapisz({
      documentId: dokument.id,
      name: nazwa.kontrolka.value.trim(),
      ...(zmieniany === '' ? {} : { templateId: zmieniany }),
      ...(przeznaczenie.value.trim() === '' ? {} : { description: przeznaczenie.value.trim() }),
      ...(kategoria.kontrolka.value.trim() === ''
        ? {}
        : { category: kategoria.kontrolka.value.trim() }),
      includeLocks: zBlokadami.checked,
    });
    if (!przyjalSie('Zapis szablonu', wynik)) return;
    if (wynik.wynik === undefined) return;
    przyjmijSzczegoly(wynik.wynik.template);
    odpowiedz.pokaz(
      `Szablon „${wynik.wynik.template.name}" zapisany jako ${wynik.wynik.template.id}. ` +
        opiszSzczegoly(wynik.wynik.template),
      true,
    );
  }

  async function ustawPoleSzablonu(usuwane: boolean): Promise<void> {
    const szablon = idSzablonu();
    if (szablon === null) return;
    if (nazwaPola.kontrolka.value.trim() === '') {
      odpowiedz.pokaz(BRAK_NAZWY_POLA, false);
      return;
    }
    const wybory = wyborWartosci.kontrolka.value
      .split(',')
      .map((czesc) => czesc.trim())
      .filter((czesc) => czesc !== '');
    const wynik = await zrodlo.ustawPole({
      templateId: szablon,
      name: nazwaPola.kontrolka.value.trim(),
      ...(etykietaPola.kontrolka.value.trim() === ''
        ? {}
        : { label: etykietaPola.kontrolka.value.trim() }),
      ...(opisPola.kontrolka.value.trim() === ''
        ? {}
        : { description: opisPola.kontrolka.value.trim() }),
      kind: rodzajPola.value as StudioTemplateFieldKind,
      required: wymagane.checked,
      ...(wartoscDomyslna.kontrolka.value.trim() === ''
        ? {}
        : { defaultValue: wartoscDomyslna.kontrolka.value.trim() }),
      ...(wybory.length === 0 ? {} : { choices: wybory }),
      ...(miejscePola.value === '' ? {} : { anchorOffset: liczba(miejscePola.value) }),
      ...(usuwane ? { remove: true } : {}),
    });
    if (!przyjalSie(usuwane ? 'Usunięcie pola szablonu' : 'Zapis pola szablonu', wynik)) return;
    if (wynik.wynik === undefined) return;
    przyjmijSzczegoly(wynik.wynik.template);
    polaSzablonu = wynik.wynik.fields;
    przerysujPola();
    odpowiedz.pokaz(
      `${usuwane ? 'Pole usunięte' : 'Pole zapisane'}: ${nazwaPola.kontrolka.value.trim()}. ` +
        `Pól szablonu po zmianie: ${polaSzablonu.length}, wymaganych ` +
        `${polaSzablonu.filter((pole) => pole.required).length}.`,
      true,
    );
  }

  async function odczytajPola(): Promise<void> {
    if (wskazany === '') {
      polaSzablonu = [];
      przerysujPola();
      return;
    }
    const wynik = await zrodlo.pola({ templateId: wskazany });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt pól szablonu', wynik.blad), false);
      return;
    }
    polaSzablonu = wynik.wynik.fields;
    przerysujPola();
  }

  async function wypelnijSzablon(): Promise<void> {
    const szablon = idSzablonu();
    if (szablon === null) return;
    const wartosci: Record<string, string> = {};
    for (const [nazwaWartosci, kontrolka] of kontrolkiWartosci) {
      wartosci[nazwaWartosci] = kontrolka.value;
    }
    const dokument = stan.dokument();
    const wynik = await zrodlo.wypelnij({
      templateId: szablon,
      windowId: stan.idOkna(),
      values: wartosci,
      ...(tytulDokumentu.kontrolka.value.trim() === ''
        ? {}
        : { title: tytulDokumentu.kontrolka.value.trim() }),
      // Dokument bieżący podaje się wyłącznie wtedy, gdy Operator go ma: brak pola
      // znaczy „załóż nowy z szablonu", a to inna czynność niż wypełnienie pól
      // w piśmie, nad którym Operator właśnie pracuje.
      ...(dokument === null ? {} : { documentId: dokument.id }),
    });
    if (!przyjalSie('Wypełnienie szablonu', wynik)) return;
    if (wynik.wynik === undefined) return;
    stan.wchlon(wynik.wynik.document);
    const brakujace = wynik.wynik.missingRequired ?? [];
    odpowiedz.pokaz(
      `Dokument z szablonu: ${wynik.wynik.document.id}. ` +
        wstawieniaOpiszBilans(wynik.wynik.balance) +
        (brakujace.length === 0
          ? ' Wszystkie pola wymagane zostały podane.'
          : ` PÓL WYMAGANYCH NIEPODANYCH: ${brakujace.length} — ${brakujace.join(', ')}. ` +
            'Dokument powstał z dziurami w tych miejscach; wypełnij je, zanim pismo wyjdzie.'),
      brakujace.length === 0,
    );
  }

  async function wniesSzablon(): Promise<void> {
    const sciezka = sciezkaWniesienia.kontrolka.value.trim();
    const plik = plikBiblioteki.kontrolka.value.trim();
    if (sciezka === '' && plik === '') {
      odpowiedz.pokaz(BRAK_PLIKU, false);
      return;
    }
    const wynik = await zrodlo.wnies({
      ...(nazwa.kontrolka.value.trim() === '' ? {} : { name: nazwa.kontrolka.value.trim() }),
      ...(formatWnoszony.value === ''
        ? {}
        : { format: formatWnoszony.value as StudioImportFormat }),
      ...(sciezka === '' ? {} : { path: sciezka }),
      ...(plik === '' ? {} : { libraryFileId: plik }),
      ...(kategoria.kontrolka.value.trim() === ''
        ? {}
        : { category: kategoria.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Wniesienie szablonu z pliku', wynik)) return;
    if (wynik.wynik === undefined) return;
    przyjmijSzczegoly(wynik.wynik.template);
    odpowiedz.pokaz(
      `Szablon „${wynik.wynik.template.name}" wniesiony jako ${wynik.wynik.template.id}. ` +
        wstawieniaOpiszWniesienie(wynik.wynik.balance) +
        ' ' +
        opiszSzczegoly(wynik.wynik.template),
      true,
    );
    void odczytajPola();
  }

  async function oddajSzablon(): Promise<void> {
    const szablon = idSzablonu();
    if (szablon === null) return;
    const wynik = await zrodlo.oddaj({
      templateId: szablon,
      ...(formatOddawany.value === ''
        ? {}
        : { format: formatOddawany.value as StudioExportFormat }),
      ...(sciezkaOddania.kontrolka.value.trim() === ''
        ? {}
        : { path: sciezkaOddania.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Oddanie szablonu do pliku', wynik)) return;
    if (wynik.wynik === undefined) return;
    odpowiedz.pokaz(wstawieniaOpiszWydanie(wynik.wynik.result), true);
  }

  async function usunSzablonWlasny(): Promise<void> {
    const szablon = idSzablonu();
    if (szablon === null) return;
    const wynik = await zrodlo.usun({ templateId: szablon });
    if (!przyjalSie('Usunięcie szablonu', wynik)) return;
    if (wynik.wynik === undefined) return;
    if (!wynik.wynik.deleted) {
      // Rdzeń odmówił — najczęściej dlatego, że szablon jest fabryczny. Pozycja
      // ZOSTAJE wskazana, bo nadal istnieje; zdjęcie jej z widoku byłoby
      // udawaniem skutku, którego nie było.
      odpowiedz.pokaz(
        `Rdzeń NIE usunął szablonu ${szablon}. Szablonu fabrycznego się nie usuwa — jest częścią ` +
          'produktu, nie wpisem Operatora. Zapisz z niego szablon własny pod nową nazwą i zmieniaj ' +
          'tamten.',
        false,
      );
      return;
    }
    const usuniety = wskazany;
    wskazany = '';
    szczegoly = null;
    polaSzablonu = [];
    przerysujPola();
    odpowiedz.pokaz(`Szablon własny ${usuniety} usunięty.`, true);
  }

  function przyjmijSzczegoly(szablon: StudioTemplateDetail): void {
    szczegoly = szablon;
    wskazany = szablon.id;
    wskazanieSzablonu.kontrolka.value = szablon.id;
    polaSzablonu = szablon.fields ?? polaSzablonu;
    przerysujPola();
  }

  function przerysujPola(): void {
    opisSzablonu.textContent =
      szczegoly === null
        ? 'Żaden szablon nie jest wskazany. Zapisz szablon z bieżącego dokumentu, wnieś go z pliku ' +
          'albo wpisz identyfikator szablonu, którego pola chcesz prowadzić.'
        : opiszSzczegoly(szczegoly);

    if (polaSzablonu.length === 0) {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Szablon nie ma ani jednego pola do wypełnienia. Wskaż miejsca, które przy każdym użyciu ' +
        'mają być podstawione: nazwa strony, data, sygnatura, adresat, kwota.';
      wykazPol.replaceChildren(puste);
      formularz.replaceChildren();
      kontrolkiWartosci.clear();
      return;
    }

    wykazPol.replaceChildren(
      ...polaSzablonu.map((pole) => {
        const wiersz = document.createElement('li');
        wiersz.dataset['pole'] = pole.name;
        wiersz.dataset['wymagane'] = pole.required ? 'tak' : 'nie';
        wiersz.textContent =
          `${pole.label} (${pole.name}) · ${opiszRodzajPola(pole.kind)}` +
          (pole.required ? ' · WYMAGANE' : ' · nieobowiązkowe') +
          (pole.defaultValue === undefined ? '' : ` · domyślnie „${pole.defaultValue}"`) +
          (pole.choices === undefined ? '' : ` · do wyboru: ${pole.choices.join(', ')}`) +
          (pole.description === undefined ? '' : ` · ${pole.description}`);
        return wiersz;
      }),
    );

    kontrolkiWartosci.clear();
    formularz.replaceChildren(
      ...polaSzablonu.map((pole) => {
        const kontrolka = document.createElement('input');
        kontrolka.type = 'text';
        kontrolka.className = 'dn-pole-kontrolka';
        kontrolka.placeholder = pole.required ? `${pole.label} (wymagane)` : pole.label;
        kontrolka.value = pole.defaultValue ?? '';
        kontrolka.setAttribute('aria-label', pole.label);
        kontrolka.dataset['pole'] = pole.name;
        kontrolkiWartosci.set(pole.name, kontrolka);
        return kontrolka;
      }),
    );
  }

  /* ── Nakładka ──────────────────────────────────────────────────────────── */

  const tresc = document.createElement('div');
  tresc.className = 'ms-szablon__tresc';
  tresc.append(
    czesc('Szablon z bieżącego dokumentu', [
      nazwa.element,
      przeznaczenie,
      kategoria.element,
      wskazanieSzablonu.element,
      etykietaBlokad,
      zapisz,
      opisSzablonu,
    ]),
    czesc('Pola do wypełnienia', [
      nazwaPola.element,
      etykietaPola.element,
      opisPola.element,
      rodzajPola,
      wartoscDomyslna.element,
      wyborWartosci.element,
      miejscePola,
      etykietaWymagania,
      ustawPole,
      usunPole,
      wykazPol,
    ]),
    czesc('Wypełnienie szablonu', [tytulDokumentu.element, formularz, wypelnij]),
    czesc('Plik Operatora', [
      sciezkaWniesienia.element,
      plikBiblioteki.element,
      formatWnoszony,
      wnies,
      sciezkaOddania.element,
      formatOddawany,
      oddaj,
      usunSzablon,
    ]),
    odpowiedz.element,
  );

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-szablon__nakladka';
  nakladka.hidden = true;
  nakladka.append(tresc);

  const wyzwalacz = przycisk('Warsztat szablonów ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'warsztat-szablonow';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera warsztat szablonów nakładką: zapis szablonu z bieżącego pisma wraz z arkuszem stylów, ' +
    'nastawami strony, nagłówkiem, stopką i blokadami wzorcowymi, pola do wypełnienia, wypełnienie, ' +
    'wniesienie z pliku dotx albo ott i oddanie do pliku. Schodzi drugim naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestaw(nakladka.hidden));

  const element = document.createElement('section');
  element.className = 'ms-szablon';
  element.setAttribute('aria-label', 'Warsztat szablonów pism');
  element.append(wyzwalacz, nakladka);
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestaw(false);
  });

  function przestaw(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    wyzwalacz.textContent = `Warsztat szablonów ${otwarta ? '▴' : '▾'}`;
    if (otwarta && wskazany !== '') void odczytajPola();
  }

  przerysujPola();

  return {
    element,
    przestawWidocznosc: () => przestaw(nakladka.hidden),
    odswiez: () => {
      if (!nakladka.hidden && wskazany !== '') void odczytajPola();
    },
    wskaz(idWskazanego) {
      wskazany = idWskazanego;
      wskazanieSzablonu.kontrolka.value = idWskazanego;
      if (nakladka.hidden) przestaw(true);
      void odczytajPola();
    },
  };
}

/** Zdanie o szablonie — co naprawdę niesie, a nie samo „zapisano". */
function opiszSzczegoly(szablon: StudioTemplateDetail): string {
  const czesci: string[] = [
    szablon.builtin ? 'szablon FABRYCZNY — usunąć się go nie da' : 'szablon własny Operatora',
    `format dokumentu ${szablon.format}`,
    `pól do wypełnienia ${(szablon.fields ?? []).length}`,
  ];
  if (szablon.category !== undefined && szablon.category !== '') {
    czesci.push(`kategoria ${szablon.category}`);
  }
  const blokady = szablon.locks ?? [];
  czesci.push(
    blokady.length === 0
      ? 'BEZ blokad wzorcowych — fragmentów, których model nie tknie, szablon nie narzuca'
      : `blokad wzorcowych ${blokady.length}, przechodzą do dokumentów z szablonu`,
  );
  czesci.push(
    szablon.form === undefined
      ? 'rdzeń nie oddał postaci wzorcowej — arkusz stylów i nastawy strony nie przyszły z tą odpowiedzią'
      : 'niesie postać wzorcową: arkusz stylów, nastawy strony, nagłówek, stopkę i tabele',
  );
  return `Szablon „${szablon.name}": ${czesci.join(' · ')}.`;
}

/** Nazwa rodzaju pola pełnym słowem. */
function opiszRodzajPola(rodzaj: StudioTemplateFieldKind): string {
  const znaleziony = RODZAJE_POL.find((pozycja) => pozycja[0] === rodzaj);
  return znaleziony === undefined ? rodzaj : znaleziony[1];
}

function liczba(wartosc: string): number {
  const odczytana = Number.parseInt(wartosc, 10);
  return Number.isFinite(odczytana) ? odczytana : 0;
}

function czesc(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-szablon__tytul';
  naglowek.textContent = tytul;
  const sekcja = document.createElement('section');
  sekcja.className = 'ms-szablon__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_DOKUMENTU =
  'Szablon powstaje z dokumentu — wczytaj pismo, doprowadź je do postaci wzorcowej i zapisz je jako ' +
  'szablon. Komenda studio.template.save bierze z dokumentu arkusz stylów i nastawy strony.';

const BRAK_NAZWY =
  'Szablon potrzebuje nazwy — jest polem obowiązkowym komendy studio.template.save i jedyną rzeczą, ' +
  'po której Operator znajdzie ten wzór w galerii.';

const BRAK_SZABLONU =
  'Ta czynność dotyczy szablonu wskazanego — wpisz jego identyfikator w polu „Szablon zmieniany" ' +
  'albo zapisz szablon z bieżącego dokumentu, a panel wskaże go sam.';

const BRAK_NAZWY_POLA =
  'Pole szablonu potrzebuje nazwy używanej w treści — to ona wiąże miejsce w piśmie z wartością ' +
  'podstawianą przy wypełnieniu.';

const BRAK_PLIKU =
  'Wniesienie szablonu potrzebuje wskazania: ścieżki widzianej przez rdzeń albo identyfikatora pliku ' +
  'Biblioteki. Klient dysku nie czyta i nie ma czego wysłać.';
