import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { WSZYSTKIE_OPERACJE, kategoriaOperacji } from './kategorie-operacji';
import type { PrzyciskMikrofonu } from './przybornik-mowa';
import type { PlywakOperacji } from './przybornik-plywak';
import type { TrybOperacji } from './przybornik-uzycie';

/** Interfejs CzynnosciWiersza niesie czynności wiersza polecenia zlecane oknu: operację, komentarz, otwarcie panelu, przypięcie, suwak i tryb wykazu. */
export interface CzynnosciWiersza {
  /** Zleca operację o wskazanym identyfikatorze wraz z poleceniem własnym. */
  naOperacje(idAkcji: string, polecenie: string): void;
  /** Zakłada komentarz przypięty do zaznaczenia. */
  naKomentarz(tresc: string): void;
  /** Otwiera pełny wykaz operacji — stały panel boczny. */
  naWiecej(): void;
  /** Przypina czynność na wierzchu pływaka albo zdejmuje przypięcie. */
  naPrzypiecie(idAkcji: string): void;
  /** Przestawia wielkość ciągłą suwakiem. */
  naSuwak(kod: string, wartosc: number): void;
  /** Przestawia tryb wykazu operacji: narzędzia ukryte albo stały panel. */
  naTryb(tryb: TrybOperacji): void;
}

/** Interfejs WierszPolecenia niesie wiersz polecenia przy kursorze wraz z jego sterowaniem: pokazaniem, ukryciem, treścią pola i zdaniem stanu. */
export interface WierszPolecenia {
  element: HTMLElement;
  /** Pokazuje wiersz przy kursorze i mówi, czego dotyczy; brak wysokości wiersza bierze odstęp domyślny. */
  pokaz(
    polozenie: { x: number; y: number },
    dlugoscZaznaczenia: number,
    wysokoscWiersza?: number,
  ): void;
  ukryj(): void;
  /** Czy wiersz jest otwarty — okno nie chowa go przy każdym odświeżeniu. */
  otwarty(): boolean;
  /** Wpisuje treść do pola polecenia — droga mikrofonu i podpowiedzi. */
  ustawTresc(tresc: string): void;
  /** Wypisuje zdanie o stanie: powodzenie albo odmowę nazwaną. */
  pokazZdanie(tresc: string, udane: boolean): void;
  /** Zwija podpowiedzi i przerywa nagranie — rozbiórka widoku. */
  zamknij(): void;
}

/**
 * Stała AKCJA_POLECENIA niesie identyfikator operacji, którą jedzie polecenie własnymi słowami, gdy Operator nie wskazał żadnej czynności z wykazu.
 */
const AKCJA_POLECENIA = 'studio.styl.rejestr';

export function utworzWierszPolecenia(
  czynnosci: CzynnosciWiersza,
  plywak: PlywakOperacji,
  mikrofon: PrzyciskMikrofonu | null,
): WierszPolecenia {
  let widoczny = false;

  const zakres = document.createElement('p');
  zakres.className = 'dn-pole-opis ms-wiersz__zakres';

  const zdanieStanu = document.createElement('p');
  zdanieStanu.className = 'dn-pole-opis ms-wiersz__zdanie';

  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole-kontrolka ms-wiersz__pole';
  pole.placeholder = 'opisz zmianę albo wpisz nazwę czynności — podpowiedzi wchodzą same';
  pole.setAttribute('aria-label', 'Polecenie dla modelu dotyczące tego fragmentu');
  pole.autocomplete = 'off';

  /** Podpowiedzi nazw czynności: wykaz płaski, bo Operator wpisuje nazwę, nie grupę katalogu. */
  const podpowiedzi = utworzMenuDrzewo({
    nastawa: 'Czynność z katalogu',
    bezUchwytu: true,
    opisTylkoPrzyWyroznionej: true,
    kierunek: 'dol',
    naWybor(klucz) {
      czynnosci.naOperacje(klucz, pole.value.trim());
      ukryj();
    },
  });
  podpowiedzi.ustaw('', drzewoPodpowiedzi());

  const zlec = document.createElement('button');
  zlec.type = 'button';
  zlec.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  zlec.textContent = 'Zleć modelowi';
  zlec.dataset['czynnosc'] = 'zlec';

  const komentarz = document.createElement('button');
  komentarz.type = 'button';
  komentarz.className = 'dn-btn dn-btn--sm dn-btn--duch';
  komentarz.textContent = 'Skomentuj';
  komentarz.dataset['czynnosc'] = 'komentarz';
  komentarz.title =
    'Zakłada komentarz redakcyjny przypięty do tego fragmentu (studio.comment.add) — ' +
    'treści nie zmienia.';

  const pasPolecenia = document.createElement('div');
  pasPolecenia.className = 'ms-wiersz__pas';
  pasPolecenia.append(pole);
  if (mikrofon !== null) pasPolecenia.append(mikrofon.element);
  pasPolecenia.append(zlec, komentarz);

  const wiecej = document.createElement('button');
  wiecej.type = 'button';
  wiecej.className = 'dn-btn dn-btn--sm dn-btn--duch';
  wiecej.textContent = 'Stały panel operacji';
  wiecej.dataset['operacja'] = 'wiecej';
  wiecej.title =
    'Przenosi ognisko do stałego panelu operacji. Panel jest trybem do wyboru — narzędzia ukryte ' +
    'zostają drogą domyślną.';
  wiecej.addEventListener('click', () => czynnosci.naWiecej());

  const element = document.createElement('div');
  element.className = 'ms-wiersz';
  element.hidden = true;
  element.append(zakres, pasPolecenia, podpowiedzi.element, plywak.element, wiecej, zdanieStanu);

  function zlecPolecenie(): void {
    if (pole.value.trim() === '') {
      zdanieStanu.textContent =
        'Napisz, co ma się zmienić, albo wybierz czynność z pływaka lub katalogu — żądanie bez ' +
        'polecenia i bez czynności nie miałoby czego zlecić.';
      return;
    }
    czynnosci.naOperacje(AKCJA_POLECENIA, pole.value);
    ukryj();
  }

  zlec.addEventListener('click', () => zlecPolecenie());
  komentarz.addEventListener('click', () => {
    if (pole.value.trim() === '') {
      zdanieStanu.textContent = 'Komentarz bez treści nie jest komentarzem — napisz go w polu wyżej.';
      return;
    }
    czynnosci.naKomentarz(pole.value);
    ukryj();
  });

  pole.addEventListener('input', () => {
    podpowiedzi.ustawFraze(pole.value.trim());
    if (pole.value.trim() === '') podpowiedzi.zwin();
    else podpowiedzi.rozwin();
  });

  pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'ArrowDown') {
      podpowiedzi.rozwin();
      podpowiedzi.przesunWyroznienie(1);
      zdarzenie.preventDefault();
      return;
    }
    if (zdarzenie.key === 'ArrowUp') {
      podpowiedzi.przesunWyroznienie(-1);
      zdarzenie.preventDefault();
      return;
    }
    if (zdarzenie.key === 'Enter') {
      // Enter najpierw wybiera podpowiedź wyróżnioną; bez wyróżnienia jedzie polecenie własnymi słowami.
      if (!podpowiedzi.wybierzWyrozniona()) zlecPolecenie();
      return;
    }
    if (zdarzenie.key === 'Escape') ukryj();
  });

  function ukryj(): void {
    widoczny = false;
    element.hidden = true;
    pole.value = '';
    podpowiedzi.zwin();
    mikrofon?.przerwij();
  }

  return {
    element,

    pokaz(polozenie, dlugoscZaznaczenia, wysokoscWiersza) {
      widoczny = true;
      element.hidden = false;
      element.style.left = `${Math.max(0, polozenie.x)}px`;
      element.style.top = `${polozenie.y + 8}px`;
      // Pływak ustawia się po odsłonięciu wiersza, bo wysokość elementu ukrytego jest zerowa.
      plywak.ustawPolozenie(polozenie, wysokoscWiersza ?? 20);
      zakres.textContent =
        dlugoscZaznaczenia > 0
          ? `Dotyczy zaznaczenia — ${dlugoscZaznaczenia} znaków. Wynik wejdzie w to miejsce jako zmiana autora „model".`
          : 'Bez zaznaczenia polecenie obejmie CAŁY dokument — zaznacz fragment, jeśli ma dotyczyć tylko jego.';
      pole.focus();
    },

    ukryj,
    otwarty: () => widoczny,

    ustawTresc(tresc) {
      pole.value = tresc;
      podpowiedzi.ustawFraze(tresc.trim());
    },

    pokazZdanie(tresc, udane) {
      zdanieStanu.textContent = tresc;
      zdanieStanu.dataset['udane'] = udane ? 'tak' : 'nie';
    },

    zamknij() {
      podpowiedzi.zwin();
      mikrofon?.przerwij();
    },
  };
}

/** Funkcja drzewoPodpowiedzi zwraca płaski wykaz podpowiedzi ze wszystkimi czynnościami katalogu operacji wraz z nazwą ich grupy. */
function drzewoPodpowiedzi(): PozycjaMenu[] {
  return WSZYSTKIE_OPERACJE.map((operacja) => ({
    rodzaj: 'wybor' as const,
    klucz: operacja.id,
    nazwa: operacja.nazwa,
    nazwaKrotka: operacja.id,
    opis: kategoriaOperacji(operacja.id)?.nazwa ?? '',
    wybrany: false,
  }));
}
