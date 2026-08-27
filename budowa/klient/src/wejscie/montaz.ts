/**
 * DROGA WEJŚCIA — montaż okien.
 *
 * Składa trzy okna ze składników i wstawia je w miejsca wskazane w dokumencie:
 *
 *     <div data-wejscie-okno="uruchomienie"></div>
 *     <div data-wejscie-okno="dostep"></div>
 *     <div data-wejscie-okno="przygotowanie"></div>
 *
 * Okno jest siatką o trzech wierszach: belka, korpus, pas działań. W wierszu
 * pasa stoją dwie rzeczy — pas przez całą szerokość i nota w kolumnie
 * tożsamości. Panele i pasy noszą ten sam `data-widok`, więc przełączają się
 * razem.
 *
 * Montaż jest jedynym miejscem, które dotyka dokumentu. Stan trzyma przebieg;
 * tutaj zostaje wyłącznie przełożenie stanu na węzły i zdarzeń na wywołania
 * przebiegu. Dzięki temu każda odsłona jest osiągalna w sprawdzianie bez
 * przeglądarki — sprawdzian prowadzi przebieg, nie okno.
 */

import { ekranDostepu, pasyDostepu } from './ekrany/dostep.ts';
import {
  PONOWIENIE,
  ekranPrzygotowania,
  pasPrzygotowania,
  paskiPrzygotowania,
  wykazPrzygotowania,
} from './ekrany/przygotowanie.ts';
import { ekranUruchomienia, pasyUruchomienia } from './ekrany/uruchomienie.ts';
import { el, naZegar, tekst } from './narzedzia.ts';
import {
  ocenHaslo,
  type Odslona,
  type Przebieg,
  type StanPrzebiegu,
  type Usterka,
} from './przebieg.ts';
import { baner } from './skladniki/baner.ts';
import {
  belkaOkna,
} from './skladniki/belka-okna.ts';
import {
  kolumnaTozsamosci,
  notaPasa,
  notaWydawcy,
} from './skladniki/kolumna-tozsamosci.ts';
import { listaEtapow } from './skladniki/lista-etapow.ts';
import { odczytajDroge, przyjmijWklejenie, zapomnijWklejenie } from './skladniki/pole-kodu.ts';
import { WARUNKI_HASLA } from './skladniki/miernik-sily.ts';

/** Ile milisekund ma takt zegara okna. Jeden zegar na całe okno. */
const TAKT_MS = 1000;

/** Powiadomienie o rzeczy, która nie ma własnej odsłony. */
export type Powiadom = (tytul: string, tresc: string) => void;

export interface NastawyMontazu {
  /** Korzeń, w którym montaż szuka miejsc na okna. */
  korzen: ParentNode;
  /** Przebieg prowadzący rozmowę z rdzeniem. */
  przebieg: Przebieg;
  /** Wersja klienta wypisywana w nocie wydawcy. */
  wersjaKlienta: string;
  /** Powiadomienie; bez wskazania idzie do obszaru komunikatów odsłony. */
  powiadom?: Powiadom;
}

/** Zamontowane okno wraz z drogą jego zdjęcia. */
export interface ZamontowaneOkno {
  /** Odłącza słuchacze i zatrzymuje zegar okna. */
  zdejmij(): void;
}

/** Trzy okna wraz z ich składem. Wykaz jest daną, nie rozgałęzieniem. */
const OKNA = {
  uruchomienie: {
    belka: 'okno.uruchamianie',
    odslona: 'uruchamianie',
    nota: null,
  },
  dostep: {
    belka: 'okno.dostep',
    odslona: 'dostep',
    nota: null,
  },
  /* Okno przygotowania nie ma belki systemowej — stoi już wewnątrz ramy
     aplikacji, a ta niesie własną. Nie ma też odsłon do przełączania. */
  przygotowanie: {
    belka: null,
    odslona: 'przygotowanie',
    nota: 'przygotowanie.nota',
  },
} as const;

type NazwaOkna = keyof typeof OKNA;

export function zamontuj(nastawy: NastawyMontazu): ZamontowaneOkno {
  const { korzen, przebieg, wersjaKlienta } = nastawy;
  const pola = [...korzen.querySelectorAll('[data-wejscie-okno]')];
  for (const pole of pola) {
    const nazwa = (pole as HTMLElement).dataset['wejscieOkno'] as NazwaOkna | undefined;
    if (nazwa === undefined || !(nazwa in OKNA)) continue;
    pole.appendChild(zbudujOkno(nazwa, przebieg.stan(), wersjaKlienta));
  }

  const powiadom: Powiadom = nastawy.powiadom ?? ((tytul, tresc) => wKomunikaty(tytul, tresc));

  zwiazZdarzenia(korzen, przebieg, powiadom);
  const odsubskrybuj = przebieg.naZmiane((stan) => odswiez(korzen, stan));
  odswiez(korzen, przebieg.stan());

  const zegar = setInterval(() => odliczaj(korzen, przebieg), TAKT_MS);
  odliczaj(korzen, przebieg);

  return {
    zdejmij() {
      clearInterval(zegar);
      odsubskrybuj();
    },
  };

  /** Powiadomienie zapasowe: informacja w obszarze komunikatów odsłony. */
  function wKomunikaty(tytul: string, tresc: string): void {
    const obszar = obszarKomunikatow(korzen, przebieg.stan().odslona);
    if (obszar === null) return;
    obszar.replaceChildren(
      baner({ rodzaj: 'informacja', ikona: 'informacja', glowa: tytul, tresc, dane: {} }),
    );
    // Głowa i treść są tu gotowym tekstem, nie kluczem — baner przyjmuje klucz,
    // więc podmieniamy oba napisy po zbudowaniu węzła.
    const wezel = obszar.firstElementChild;
    const glowa = wezel?.querySelector('b');
    if (glowa !== null && glowa !== undefined) glowa.textContent = tytul;
    const trescWezla = wezel?.querySelector('.dn-alert-tresc');
    if (trescWezla !== null && trescWezla !== undefined) {
      trescWezla.lastChild!.textContent = tresc;
    }
  }
}

/* ── Budowa ──────────────────────────────────────────────────────────────── */

function zbudujOkno(nazwa: NazwaOkna, stan: StanPrzebiegu, wersjaKlienta: string): HTMLElement {
  const opis = OKNA[nazwa];
  const dzieci: (HTMLElement | null)[] = [];
  if (opis.belka !== null) dzieci.push(belkaOkna({ tytul: opis.belka }));
  dzieci.push(kolumnaTozsamosci({ odslona: opis.odslona }));

  if (nazwa === 'uruchomienie') {
    dzieci.push(...ekranUruchomienia(), ...pasyUruchomienia());
  } else if (nazwa === 'dostep') {
    dzieci.push(...ekranDostepu(), ...pasyDostepu());
  } else {
    dzieci.push(ekranPrzygotowania(stan.przygotowanie), pasPrzygotowania());
  }

  dzieci.push(opis.nota === null ? notaWydawcy(wersjaKlienta) : notaPasa(opis.nota));
  return el(
    'div',
    { klasa: `we-okno${opis.belka === null ? '' : ' we-okno--z-belka'}` },
    dzieci,
  );
}

/* ── Odświeżenie ─────────────────────────────────────────────────────────── */

function odswiez(korzen: ParentNode, stan: StanPrzebiegu): void {
  przelaczOdslony(korzen, stan.odslona);
  odswiezEtapyLaczenia(korzen, stan);
  odswiezZwloke(korzen, stan);
  odswiezAdresy(korzen, stan);
  odswiezPrzygotowanie(korzen, stan);
  odswiezUsterki(korzen, stan);
  odswiezOczekiwanie(korzen, stan);
}

/** Panel i pas noszą ten sam `data-widok`, więc przełączają się razem. */
function przelaczOdslony(korzen: ParentNode, odslona: Odslona): void {
  for (const wezel of korzen.querySelectorAll('[data-widok]')) {
    const czynny = (wezel as HTMLElement).dataset['widok'] === odslona;
    (wezel as HTMLElement).dataset['widokAktywny'] = czynny ? 'tak' : 'nie';
  }
  for (const zakladka of korzen.querySelectorAll('[role="tab"][data-idz]')) {
    const panel = zakladka.closest('[data-pigulka]') as HTMLElement | null;
    const wybrana = panel?.dataset['pigulka'];
    zakladka.setAttribute(
      'aria-selected',
      (zakladka as HTMLElement).dataset['idz'] === wybrana ? 'true' : 'false',
    );
  }
}

function odswiezEtapyLaczenia(korzen: ParentNode, stan: StanPrzebiegu): void {
  for (const gniazdo of korzen.querySelectorAll('[data-etapy-laczenia]')) {
    if ((gniazdo as HTMLElement).dataset['etapyLaczenia'] !== stan.odslona) continue;
    gniazdo.replaceChildren(
      listaEtapow({
        nazwy: 'uruchomienie.etapy',
        miary: 'uruchomienie.stany',
        obszar: 'uruchomienie.laczenie.tytul',
        stany: stan.etapyLaczenia,
        wartosciMiar: stan.miaryLaczenia.map((klucz) => ({ klucz })),
      }),
    );
  }
}

/**
 * Licznik zwłoki dostaje wartość ZMIERZONĄ przez przebieg — tyle, ile rdzeń
 * kazał czekać przy próbie poprzedniej. Odmowa rdzenia tej wartości nie niesie,
 * więc nie ma jej skąd odczytać; wpisanie tu stałej byłoby obietnicą czasu,
 * którego nikt nie mierzył.
 */
function odswiezZwloke(korzen: ParentNode, stan: StanPrzebiegu): void {
  if (stan.zwlokaS <= 0) return;
  const panel = korzen.querySelector(`[data-widok="${stan.odslona}"]`);
  const licznik = panel?.querySelector('[data-odliczanie]') as HTMLElement | null;
  if (licznik === null || licznik === undefined) return;
  // Nastawiamy licznik wyłącznie wtedy, gdy stoi na zerze, czyli przy wejściu
  // w odsłonę. Nastawianie go przy każdej zmianie stanu zamroziłoby odliczanie.
  if (Number.parseInt(licznik.dataset['odliczanie'] ?? '0', 10) > 0) return;
  licznik.dataset['odliczanie'] = String(stan.zwlokaS);
}

function odswiezAdresy(korzen: ParentNode, stan: StanPrzebiegu): void {
  for (const wezel of korzen.querySelectorAll('[data-adres]')) {
    wezel.textContent = stan.adres;
  }
  const glowa = korzen.querySelector(
    '[data-widok="konto-bez-potwierdzenia"] .dn-alert-tresc b',
  );
  if (glowa !== null) {
    glowa.textContent = tekst('dostep.kontoBezPotwierdzenia.baner.glowa', { adres: stan.adres });
  }
}

function odswiezPrzygotowanie(korzen: ParentNode, stan: StanPrzebiegu): void {
  const wykaz = korzen.querySelector('[data-etapy-przygotowania]');
  if (wykaz !== null) wykaz.replaceChildren(wykazPrzygotowania(stan.przygotowanie));
  const pasek = korzen.querySelector('[data-postep-przygotowania]');
  if (pasek !== null) pasek.replaceChildren(paskiPrzygotowania(stan.przygotowanie));
  for (const miara of korzen.querySelectorAll('.dn-postep-wartosc[data-wartosc]')) {
    (miara as HTMLElement).style.width = `${(miara as HTMLElement).dataset['wartosc'] ?? '0'}%`;
  }
  // Ponowienie stoi wyłącznie przy etapie nieudanym: przy przebiegu udanym nie
  // ma czego ponawiać, a kontrolka bez skutku jest gorsza od jej braku.
  const ponowienie = korzen.querySelector(`#${PONOWIENIE}`) as HTMLElement | null;
  if (ponowienie !== null) ponowienie.hidden = !stan.przygotowanie.stany.includes('blad');
}

function obszarKomunikatow(korzen: ParentNode, odslona: Odslona): Element | null {
  return korzen.querySelector(`[data-komunikaty="${odslona}"]`);
}

/**
 * Usterki zbierają się w JEDEN komunikat. Kilka osobnych banerów zepchnęłoby
 * formularz poza okno, a użytkownik i tak czyta je jako jedną listę tego, co
 * ma poprawić. Pole, którego usterka dotyczy, niesie `aria-invalid` — czytnik
 * ekranu dowiaduje się tego samego co oko.
 */
function odswiezUsterki(korzen: ParentNode, stan: StanPrzebiegu): void {
  for (const obszar of korzen.querySelectorAll('[data-komunikaty]')) obszar.replaceChildren();
  for (const pole of korzen.querySelectorAll('[aria-invalid="true"][data-usterka-pola]')) {
    pole.removeAttribute('aria-invalid');
  }
  if (stan.usterki.length === 0) return;
  const obszar = obszarKomunikatow(korzen, stan.odslona);
  if (obszar === null) return;

  const odmowa = stan.usterki.find((u) => u.odRdzenia !== undefined);
  if (odmowa?.odRdzenia !== undefined) {
    // Rdzeń wie o powodzie odmowy więcej niż okno, więc treścią jest to, co
    // powiedział rdzeń; katalog daje wyłącznie głowę komunikatu.
    const wezel = baner({
      rodzaj: 'blad',
      ikona: 'ostrzezenie',
      glowa: 'usterki.odmowaRdzenia.glowa',
      tresc: 'usterki.odmowaRdzenia.glowa',
    });
    const tresc = wezel.querySelector('.dn-alert-tresc');
    if (tresc?.lastChild != null) tresc.lastChild.textContent = odmowa.odRdzenia.message;
    obszar.replaceChildren(wezel);
    return;
  }

  const klucze = stan.usterki.map((u) => u.klucz).filter((k): k is string => k !== undefined);
  if (klucze.length === 0) return;
  if (klucze.length === 1) {
    obszar.replaceChildren(
      baner({
        rodzaj: 'blad',
        ikona: 'ostrzezenie',
        glowa: kluczUsterki(klucze[0]!, 'glowa'),
        tresc: kluczUsterki(klucze[0]!, 'tresc'),
      }),
    );
  } else {
    /* Przy kilku usterkach wykaz niesie same rozpoznania, bez wskazówek: co
       zrobić, mówi zaznaczone pole, a wskazówki powtórzone kilka razy
       wypchnęłyby formularz poza okno. */
    const wezel = baner({
      rodzaj: 'blad',
      ikona: 'ostrzezenie',
      glowa: 'usterki.wiele',
      trescWezly: [
        el(
          'ul',
          { klasa: 'dn-alert-lista' },
          klucze.map((k) => el('li', { tekst: tekst(kluczUsterki(k, 'glowa')) })),
        ),
      ],
    });
    obszar.replaceChildren(wezel);
  }
  oznaczPolaUsterek(korzen, stan.odslona, klucze);
}

/** Klucz katalogu dla rozpoznania usterki; nazwy z kropką są już pełne. */
function kluczUsterki(klucz: string, czesc: 'glowa' | 'tresc'): string {
  if (klucz.includes('.')) return `${klucz}.${czesc}`;
  return `usterki.${zNazwyKreskowej(klucz)}.${czesc}`;
}

/** Nazwa kreskowa na nazwę katalogu: `haslo-slabe` → `hasloSlabe`. */
function zNazwyKreskowej(nazwa: string): string {
  return nazwa.replace(/-([a-z])/g, (_, znak: string) => znak.toUpperCase());
}

/** Pola, których dotyczy każde rozpoznanie — kontrakt z formularzem. */
const POLA_USTEREK: Record<string, string[]> = {
  brakLoginu: ['log-login', 'blad-login', 'rej-login'],
  brakHasla: ['log-haslo', 'blad-haslo', 'rej-haslo', 'odz-haslo'],
  brakDanych: ['log-login', 'log-haslo', 'blad-login', 'blad-haslo'],
  brakAdresu: ['rej-email', 'odz-email'],
  brakDrogi: [],
  'email-bledny': ['rej-email', 'odz-email'],
  'haslo-slabe': ['rej-haslo', 'odz-haslo'],
  'hasla-rozne': ['rej-haslo', 'rej-haslo-2', 'odz-haslo', 'odz-haslo-2'],
};

function oznaczPolaUsterek(korzen: ParentNode, odslona: Odslona, klucze: string[]): void {
  const panel = korzen.querySelector(`[data-widok="${odslona}"]`);
  if (panel === null) return;
  let pierwsze: HTMLElement | null = null;
  for (const klucz of klucze) {
    for (const id of POLA_USTEREK[klucz] ?? []) {
      const pole = panel.querySelector(`#${id}`) as HTMLElement | null;
      if (pole === null) continue;
      pole.setAttribute('aria-invalid', 'true');
      pole.dataset['usterkaPola'] = '';
      pierwsze ??= pole;
    }
  }
  pierwsze?.focus();
}

/**
 * Zero blokad: kontrolka niegotowa nie jest wyłączana, tylko oznaczona jako
 * zajęta. Przycisk wygaszony nie mówi, dlaczego nie działa.
 */
function odswiezOczekiwanie(korzen: ParentNode, stan: StanPrzebiegu): void {
  for (const pas of korzen.querySelectorAll('.we-pas .dn-btn--sygnal')) {
    if (stan.wToku) pas.setAttribute('aria-busy', 'true');
    else pas.removeAttribute('aria-busy');
  }
}

/* ── Zdarzenia ───────────────────────────────────────────────────────────── */

function zwiazZdarzenia(korzen: ParentNode, przebieg: Przebieg, powiadom: Powiadom): void {
  const dokument = korzen as ParentNode & { addEventListener?: typeof document.addEventListener };
  dokument.addEventListener?.('click', (zdarzenie: Event) => {
    const cel = (zdarzenie.target as Element | null)?.closest(
      '[data-idz], [data-czynnosc], [data-komunikat], [data-odsloniecie], [data-wklej-kod]',
    ) as HTMLElement | null;
    if (cel === null) return;

    if (cel.dataset['odsloniecie'] !== undefined) {
      zdarzenie.preventDefault();
      odslonHaslo(korzen, cel);
      return;
    }
    if (cel.dataset['wklejKod'] !== undefined) {
      zdarzenie.preventDefault();
      void wklejDroge(korzen, cel.dataset['wklejKod']!, powiadom);
      return;
    }
    if (cel.dataset['komunikat'] !== undefined) {
      zdarzenie.preventDefault();
      powiadom(cel.dataset['komunikatTytul'] ?? '', cel.dataset['komunikat']);
      return;
    }
    if (cel.dataset['idz'] !== undefined) {
      zdarzenie.preventDefault();
      przebieg.przejdzDo(cel.dataset['idz'] as never);
      return;
    }
    zdarzenie.preventDefault();
    void wykonaj(korzen, przebieg, powiadom, cel.dataset['czynnosc']!);
  });

  dokument.addEventListener?.('input', (zdarzenie: Event) => {
    const cel = zdarzenie.target as HTMLElement | null;
    if (cel === null) return;
    if (cel.classList.contains('au-kod-pole')) obsluzZnakDrogi(korzen, cel as HTMLInputElement);
    odswiezMierniki(korzen);
  });

  dokument.addEventListener?.('paste', (zdarzenie: Event) => {
    const cel = zdarzenie.target as HTMLElement | null;
    if (cel === null || !cel.classList.contains('au-kod-pole')) return;
    const zestaw = cel.closest('[data-kod-grupa]') as HTMLElement | null;
    if (zestaw === null) return;
    zdarzenie.preventDefault();
    const dane = (zdarzenie as ClipboardEvent).clipboardData?.getData('text') ?? '';
    przyjmijWklejenie(zestaw, dane);
  });

  odswiezMierniki(korzen);
}

/** Odsłonięcie hasła zmienia typ pola i własną etykietę. */
function odslonHaslo(korzen: ParentNode, kontrolka: HTMLElement): void {
  const pole = korzen.querySelector(`#${kontrolka.dataset['odsloniecie']}`) as HTMLInputElement | null;
  if (pole === null) return;
  const bylo = pole.type === 'text';
  pole.type = bylo ? 'password' : 'text';
  kontrolka.setAttribute('aria-pressed', bylo ? 'false' : 'true');
  kontrolka.setAttribute('aria-label', tekst(bylo ? 'dostep.haslo.pokaz' : 'dostep.haslo.ukryj'));
}

async function wklejDroge(korzen: ParentNode, grupa: string, powiadom: Powiadom): Promise<void> {
  const zestaw = korzen.querySelector(`[data-kod-grupa="${grupa}"]`) as HTMLElement | null;
  if (zestaw === null) return;
  const schowek = globalThis.navigator?.clipboard;
  if (schowek === undefined || typeof schowek.readText !== 'function') {
    powiadom(tekst('komunikaty.schowek.tytul'), tekst('komunikaty.schowek.tresc'));
    return;
  }
  try {
    przyjmijWklejenie(zestaw, await schowek.readText());
  } catch {
    powiadom(tekst('komunikaty.schowek.tytul'), tekst('komunikaty.schowek.tresc'));
  }
}

/** Przejście między polami drogi potwierdzenia; wpisanie znosi wklejenie. */
function obsluzZnakDrogi(korzen: ParentNode, pole: HTMLInputElement): void {
  const zestaw = pole.closest('[data-kod-grupa]') as HTMLElement | null;
  if (zestaw === null) return;
  zapomnijWklejenie(zestaw);
  pole.value = pole.value.slice(0, 1);
  const pola = [...zestaw.querySelectorAll('input')];
  const numer = pola.indexOf(pole);
  if (pole.value.length > 0) pola[numer + 1]?.focus();
  void korzen;
}

/** Miernik siły ocenia tą samą regułą, którą przebieg sprawdza hasło. */
function odswiezMierniki(korzen: ParentNode): void {
  for (const blok of korzen.querySelectorAll('[data-sila-dla]')) {
    const miernik = blok as HTMLElement;
    const pole = korzen.querySelector(`#${miernik.dataset['silaDla']}`) as HTMLInputElement | null;
    if (pole === null) continue;
    const wynik = ocenHaslo(pole.value);
    let stopien = 0;
    for (const warunek of WARUNKI_HASLA) {
      const wezel = miernik.querySelector(`[data-warunek="${warunek}"]`);
      const spelniony = wynik[warunek] === true;
      wezel?.setAttribute('data-spelniony', spelniony ? 'tak' : 'nie');
      if (spelniony) stopien += 1;
    }
    miernik.dataset['stopien'] = String(stopien);
    const opis = miernik.querySelector('[data-sila-opis]');
    if (opis !== null) {
      opis.textContent =
        pole.value.length === 0
          ? tekst('dostep.sila.puste')
          : tekst(`dostep.sila.stopnie.${stopien}`);
    }
  }
}

/** Wartość pola formularza; brak pola daje pusty łańcuch, nie wyjątek. */
function wartosc(korzen: ParentNode, id: string): string {
  return (korzen.querySelector(`#${id}`) as HTMLInputElement | null)?.value ?? '';
}

function zaznaczone(korzen: ParentNode, id: string): boolean {
  return (korzen.querySelector(`#${id}`) as HTMLInputElement | null)?.checked === true;
}

async function wykonaj(
  korzen: ParentNode,
  przebieg: Przebieg,
  powiadom: Powiadom,
  czynnosc: string,
): Promise<void> {
  switch (czynnosc) {
    case 'ponow-polaczenie':
      przebieg.ponow();
      return;
    case 'zaloguj':
      return przebieg.zaloguj({
        login: wartosc(korzen, 'log-login'),
        haslo: wartosc(korzen, 'log-haslo'),
        niewylogowuj: zaznaczone(korzen, 'log-sesja'),
      });
    case 'zaloguj-ponownie':
      return przebieg.zaloguj({
        login: wartosc(korzen, 'blad-login'),
        haslo: wartosc(korzen, 'blad-haslo'),
      });
    case 'zarejestruj':
      return przebieg.zarejestruj({
        login: wartosc(korzen, 'rej-login'),
        email: wartosc(korzen, 'rej-email'),
        haslo: wartosc(korzen, 'rej-haslo'),
        hasloPowtorzone: wartosc(korzen, 'rej-haslo-2'),
        niewylogowuj: zaznaczone(korzen, 'rej-sesja'),
      });
    case 'potwierdz-adres':
      return przebieg.potwierdzAdres({
        droga: odczytajDroge(korzen, 'kod'),
        niewylogowuj: zaznaczone(korzen, 'rej-sesja'),
      });
    case 'wejdz-bez-potwierdzenia':
      return przebieg.wejdzBezPotwierdzenia();
    case 'popros-o-odzyskanie':
      return przebieg.poprosOOdzyskanie(wartosc(korzen, 'odz-email'));
    case 'przejdz-do-hasla':
      przebieg.przejdzDo('odzyskiwanie-haslo');
      return;
    case 'ustaw-nowe-haslo':
      return przebieg.ustawNoweHaslo({
        droga: odczytajDroge(korzen, 'odzyskiwanie'),
        haslo: wartosc(korzen, 'odz-haslo'),
        hasloPowtorzone: wartosc(korzen, 'odz-haslo-2'),
      });
    case 'ponow-droge':
      powiadom(tekst('komunikaty.kodPonowiony.tytul'), tekst('komunikaty.kodPonowiony.tresc'));
      return przebieg.poprosOOdzyskanie(przebieg.stan().adres);
    case 'ponow-przygotowanie':
      return przebieg.przygotujSrodowisko();
    case 'przerwij-i-wyloguj':
      przebieg.przejdzDo('logowanie');
      return;
    default:
      return;
  }
}

/* ── Zegar okna ──────────────────────────────────────────────────────────── */

/**
 * Każdy węzeł z `data-odliczanie` niesie czas w sekundach i sam się wypisuje.
 * Napis, który nie ubywa, jest gorszy od braku napisu: obiecuje odmierzanie,
 * którego nie ma. Jeden zegar na całe okno — nie kilkanaście osobnych.
 *
 * Wypisywane są WSZYSTKIE liczniki, także w odsłonach ukrytych — inaczej
 * licznik byłby pusty przez sekundę po pokazaniu odsłony. Ubywa natomiast
 * tylko licznik widoczny: czas, który schodzi za plecami, doprowadza do tego,
 * że użytkownik zastaje zero, choć odsłonę zobaczył przed chwilą.
 */
function odliczaj(korzen: ParentNode, przebieg: Przebieg): void {
  const odslona = przebieg.stan().odslona;
  for (const wezel of korzen.querySelectorAll('[data-odliczanie]')) {
    const licznik = wezel as HTMLElement;
    const zostalo = Number.parseInt(licznik.dataset['odliczanie'] ?? '', 10);
    if (Number.isNaN(zostalo)) continue;
    const napis =
      licznik.dataset['odliczaniePostac'] === 'sekundy' ? String(zostalo) : naZegar(zostalo);
    if (licznik.textContent !== napis) licznik.textContent = napis;
    const panel = licznik.closest('[data-widok]') as HTMLElement | null;
    if (panel !== null && panel.dataset['widok'] !== odslona) continue;
    if (zostalo <= 0) {
      poZerze(licznik, przebieg, odslona);
      continue;
    }
    licznik.dataset['odliczanie'] = String(zostalo - 1);
  }
}

/**
 * Co dzieje się po dojściu do zera, rozstrzyga odsłona, w której licznik stoi.
 * Wstrzymanie mija — okno wraca tam, skąd je wstrzymano. Serwer nie odpowiadał
 * — okno ponawia próbę samo, tak jak zapowiada baner. Ważność drogi
 * potwierdzenia wygasa; okno licznik zatrzymuje i niczego nie udaje, bo rdzeń
 * nie ma komendy wydającej drogę potwierdzenia adresu po raz drugi.
 */
function poZerze(licznik: HTMLElement, przebieg: Przebieg, odslona: Odslona): void {
  if (odslona === 'w-blad') {
    licznik.dataset['odliczanie'] = licznik.dataset['odliczaniePoczatek'] ?? '0';
    przebieg.ponow();
    return;
  }
  if (odslona === 'logowanie-wstrzymane' || odslona === 'odzyskiwanie-wstrzymane') {
    przebieg.zdejmijWstrzymanie();
  }
}

/** Usterki wystawione sprawdzianowi montażu; okno ich nie tworzy samo. */
export type { Usterka };
