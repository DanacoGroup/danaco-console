import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { WSZYSTKIE_OPERACJE, kategoriaOperacji } from './kategorie-operacji';
import type { PrzyciskMikrofonu } from './przybornik-mowa';
import type { PlywakOperacji } from './przybornik-plywak';
import type { TrybOperacji } from './przybornik-uzycie';

/**
 * Wiersz polecenia przy kursorze — droga trzecia do czynności modelu.
 *
 * ── Cztery drogi do tej samej czynności ─────────────────────────────────────
 * Rozstrzygnięcie Właściciela o narzędziach ukrytych wymienia je wprost:
 * pływak przy zaznaczeniu, uchwyt pełnego katalogu, ten wiersz polecenia
 * i skrót „/" w oknie rozmowy sesji. Wykaz czynności jest przy tym JEDEN —
 * `kategorie-operacji.ts` — więc żadna droga nie zna czynności, których nie
 * znają pozostałe.
 *
 * ── Dlaczego przy kursorze, a nie w panelu ──────────────────────────────────
 * Model pracuje tam, gdzie stoi Operator. Pole polecenia otwiera się w miejscu
 * zaznaczenia, wynik wchodzi w to samo miejsce jako zmiana oznaczona autorem,
 * a Operator nie przenosi się do żadnego okna. Wiersz nosi przy tym cały pływak
 * narzędzi ukrytych: czynności najczęstsze, uchwyt katalogu i suwaki wielkości
 * ciągłych.
 *
 * ── Podpowiadanie nazw czynności ────────────────────────────────────────────
 * Operator pisze polecenie własnymi słowami, a wiersz podpowiada mu nazwy
 * czynności z katalogu — mechanizmem `komponenty/menu-drzewo.ts` w trybie bez
 * uchwytu, tym samym, którym obsadzony jest skrót „/" w oknie rozmowy. Strzałki
 * chodzą po podpowiedziach, Enter wybiera wyróżnioną, a gdy żadnej nie ma —
 * zleca polecenie własnymi słowami. Ognisko zostaje w polu, więc pisanie nie
 * jest przerywane.
 *
 * ── Czym jedzie polecenie ───────────────────────────────────────────────────
 * Polem `params` żądania `studio.contextual.op`, wraz z nastawami suwaków.
 * Rdzeń dokłada `params` do treści polecenia dla modelu, więc słowa Operatora
 * dojeżdżają tą samą drogą, którą jedzie zakres zaznaczenia. Identyfikator
 * akcji zostaje przy tym prawdziwy — pochodzi z wykazu, nie ze zdania Operatora,
 * bo `actionId` jest pozycją rejestru, nie polem na wypowiedź.
 */

/** Czynności wiersza polecenia zlecane oknu. */
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

/** Wiersz polecenia wraz z jego sterowaniem. */
export interface WierszPolecenia {
  element: HTMLElement;
  /**
   * Pokazuje wiersz przy kursorze i mówi, czego dotyczy.
   *
   * `wysokoscWiersza` służy ustawieniu pływaka pod zaznaczeniem, gdy nad nim nie
   * ma miejsca; jej brak bierze wysokość wiersza treści za odstęp domyślny.
   */
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
 * Operacja, którą jedzie polecenie własnymi słowami, gdy Operator nie wskazał
 * żadnej z wykazu.
 *
 * Nie jest to „operacja dowolna": `actionId` musi być pozycją wykazu, więc
 * polecenie własne jedzie przez przepisanie fragmentu, a treść polecenia
 * rozstrzyga, co model ma zrobić.
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

  /**
   * Podpowiedzi nazw czynności.
   *
   * Wykaz jest płaski, bo Operator wpisuje nazwę, a nie wędruje po grupach —
   * grupy są w katalogu pod uchwytem. Opis rysuje się przy pozycji wyróżnionej,
   * żeby 28 opisów naraz nie zasłoniło samych nazw.
   */
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
      // Enter najpierw wybiera podpowiedź wyróżnioną; brak wyróżnienia znaczy,
      // że Operator pisze polecenie własnymi słowami, i wtedy jedzie ono.
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
      // Pływak ustawia się dopiero po odsłonięciu wiersza: wysokość elementu
      // ukrytego jest zerowa, więc rachunek „nad czy pod" liczony przed
      // odsłonięciem postawiłby go na zaznaczeniu.
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

/** Płaski wykaz podpowiedzi — wszystkie czynności katalogu wraz z ich grupą. */
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
