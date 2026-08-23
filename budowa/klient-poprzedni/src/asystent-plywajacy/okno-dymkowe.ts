import { elementIkony } from '../ikony/ikony';
import { liczbaOkienRozmowy } from '../okno-komunikacji/profil-modulu';
import { profilModulu } from '../okno-komunikacji/rejestr-profilow';
import { czyGlosDziala, odmowMowy, powodOdmowyGlosu } from './dostepnosc-mowy';
import { zdanieOPracy } from './favikon-plywajacy';
import { oknoGotowe, zdanieOStanieOkna, type StanDymka, type Wypowiedz } from './stan-dymka';

/**
 * Okno dymkowe Asystenta — rozmowa prowadzona przy pływającym favikonie,
 * zamiast osobnego okna czatu.
 *
 * Postać bierze się z profilu, nie z tego pliku. `profilModulu('assistant')`
 * niesie `postacRozmowy: 'dymek-glosowy'`, `granicaOkien: 1` i
 * `pamiecSesyjna: true`. Dymek pyta o liczbę okien funkcją `liczbaOkienRozmowy`,
 * nie polem `granicaOkien`: pole niesie granicę, funkcja niesie prawo do
 * otwarcia. Profil przestawiony na `postacRozmowy: 'brak'` daje zdanie
 * o rozbieżności, a nie zniknięcie dymka.
 *
 * Pasek pod nagłówkiem mówi o braku kanału głosowego od chwili otwarcia, więc
 * Operator dowiaduje się o nim, zanim sięgnie po mikrofon. Przycisk mikrofonu
 * nie jest wygaszany i odmawia z powodem, który wymienia brakujące ogniwa
 * z nazwy. Nagrywania do pamięci tu nie ma: nagranie, którego nie ma dokąd
 * wysłać, byłoby atrapą mikrofonu.
 *
 * Wiersz nad polem wypowiedzi niesie zdanie o oknie modułu — osobne dla odmowy
 * `window.list`, osobne dla braku okna. Dopóki polecenie nie ma dokąd pojechać,
 * Operator widzi dlaczego.
 */

/**
 * Etykieta wiersza rozmowy — czytnik ekranu ma wiedzieć, kto mówi.
 *
 * `posuniecie` nie jest mówcą, tylko rodzajem wiersza: asystent nic nie
 * powiedział, tylko coś w aplikacji zrobił. Etykieta mówi o czynności.
 */
const ETYKIETY: Readonly<Record<Wypowiedz['rodzaj'], string>> = {
  operator: 'Operator',
  asystent: 'Asystent',
  rdzen: 'Rdzeń',
  odmowa: 'Odmowa',
  posuniecie: 'Posunięcie w aplikacji',
};

export interface OknoDymkowe {
  element: HTMLElement;
  /** Odsłania dymek i ustawia ognisko na polu wypowiedzi. */
  pokaz(): void;
  /** Chowa dymek; rozmowa zostaje (profil: pamięć sesyjna). */
  schowaj(): void;
  widoczny(): boolean;
  rozlacz(): void;
}

export function utworzOknoDymkowe(stan: StanDymka, naZamkniecie: () => void): OknoDymkowe {
  const profil = profilModulu('assistant');
  const ileOkien = liczbaOkienRozmowy(profil);

  const element = document.createElement('section');
  element.className = 'ap-dymek';
  element.hidden = true;
  element.setAttribute('aria-label', `${profil.nazwa} — okno dymkowe rozmowy`);

  // ── Nagłówek ────────────────────────────────────────────────────────────
  const glowa = document.createElement('header');
  glowa.className = 'ap-dymek__glowa';

  const tytul = document.createElement('div');
  tytul.className = 'ap-dymek__tytul';
  const nazwa = document.createElement('strong');
  nazwa.textContent = profil.nazwa;
  const podtytul = document.createElement('span');
  podtytul.className = 'ap-dymek__podtytul';
  podtytul.textContent = profil.przeznaczenie;
  tytul.append(nazwa, podtytul);

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn-ikona ap-dymek__zamknij';
  zamknij.title = 'Zwiń dymek — rozmowa zostaje';
  zamknij.setAttribute('aria-label', 'Zwiń dymek Asystenta');
  zamknij.append(elementIkony('zamknij', { rozmiar: 16 }));
  zamknij.addEventListener('click', naZamkniecie);

  glowa.append(tytul, zamknij);

  // ── Wskaźnik pracy asystenta ────────────────────────────────────────────
  // Ta sama prawda, którą niesie favikon i pas dolny, rozwinięta w zdanie.
  // Wiersz stoi zawsze: „nie prowadzi zlecenia" jest odpowiedzią, nie pustką,
  // a wiersz pojawiający się i znikający przeskakiwałby treścią.
  const wskaznikPracy = document.createElement('p');
  wskaznikPracy.className = 'ap-dymek__praca';
  wskaznikPracy.setAttribute('role', 'status');

  // ── Pas prawdy o głosie ─────────────────────────────────────────────────
  // Pas stoi zawsze; gdy kanał głosowy jest dostępny, mówi o dostępności.
  const pasGlosu = document.createElement('button');
  pasGlosu.type = 'button';
  pasGlosu.className = 'ap-dymek__glos';
  pasGlosu.dataset.dziala = String(czyGlosDziala());
  pasGlosu.textContent = czyGlosDziala()
    ? 'Kanał głosowy: dostępny.'
    : 'Kanał głosowy: NIE MA GO. Rozmowa idzie tekstem — naciśnij, żeby poznać powód.';
  pasGlosu.addEventListener('click', odmowMowy);

  // ── Postać rozmowy wzięta z profilu ─────────────────────────────────────
  // Zdanie powstaje tylko wtedy, gdy profil mówi co innego niż dymek głosowy.
  // Milczące zbudowanie dymka wbrew profilowi byłoby drugą prawdą o module.
  const rozbieznosc = zdanieORozbieznosci(profil.postacRozmowy, ileOkien);
  const pasProfilu = document.createElement('p');
  pasProfilu.className = 'ap-dymek__profil';
  pasProfilu.hidden = rozbieznosc === '';
  pasProfilu.textContent = rozbieznosc;

  // ── Historia rozmowy ────────────────────────────────────────────────────
  const historia = document.createElement('div');
  historia.className = 'ap-dymek__historia';
  historia.setAttribute('role', 'log');
  historia.setAttribute('aria-live', 'polite');
  historia.setAttribute('aria-label', 'Przebieg rozmowy z Asystentem');

  // ── Stan okna modułu ────────────────────────────────────────────────────
  const stanOkna = document.createElement('p');
  stanOkna.className = 'ap-dymek__okno';
  stanOkna.setAttribute('role', 'status');

  // ── Pole wypowiedzi ─────────────────────────────────────────────────────
  const dol = document.createElement('form');
  dol.className = 'ap-dymek__dol';

  const pole = document.createElement('textarea');
  pole.className = 'ap-dymek__pole';
  pole.rows = 2;
  pole.placeholder = 'Napisz polecenie — pojedzie komendą assistant.voice.command';
  pole.setAttribute('aria-label', 'Treść polecenia dla Asystenta');

  const mikrofon = document.createElement('button');
  mikrofon.type = 'button';
  mikrofon.className = 'dn-btn-ikona ap-dymek__mikrofon';
  mikrofon.dataset.dziala = String(czyGlosDziala());
  mikrofon.title = powodOdmowyGlosu();
  mikrofon.setAttribute('aria-label', 'Powiedz polecenie głosem — kanał głosowy niezbudowany');
  mikrofon.append(elementIkony('mikrofon', { rozmiar: 20 }));

  const wyslij = document.createElement('button');
  wyslij.type = 'submit';
  wyslij.className = 'dn-btn dn-btn--sygnal dn-btn--sm ap-dymek__wyslij';
  wyslij.textContent = 'Wyślij';

  dol.append(pole, mikrofon, wyslij);

  element.append(glowa, wskaznikPracy, pasGlosu, pasProfilu, historia, stanOkna, dol);

  // ── Zachowania ──────────────────────────────────────────────────────────

  // Ślad w historii idzie przy pierwszym naciśnięciu. Powód jest długi i za
  // każdym razem ten sam, więc powtórzony zasypałby rozmowę. Powiadomienie
  // odpowiada na każde naciśnięcie.
  let odnotowanoOdmoweGlosu = false;
  mikrofon.addEventListener('click', () => {
    odmowMowy();
    if (!odnotowanoOdmoweGlosu) {
      odnotowanoOdmoweGlosu = true;
      stan.odnotujOdmowe(powodOdmowyGlosu());
    }
    pole.focus();
  });

  dol.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    const tresc = pole.value;
    pole.value = '';
    void stan.wyslij(tresc);
  });

  // Enter wysyła, Shift+Enter łamie wiersz — jak w oknie komunikacji.
  pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    dol.requestSubmit();
  });

  /** Przerysowuje historię i wiersze stanu. */
  function odswiez(): void {
    historia.replaceChildren(...stan.wypowiedzi().map(wierszWypowiedzi));
    historia.scrollTop = historia.scrollHeight;

    // Zdanie o pracy składa `favikon-plywajacy.ts` — jedno źródło napisu dla
    // przycisku i dla dymka, żeby dwa miejsca nie nazwały tego samego inaczej.
    const praca = stan.praca();
    wskaznikPracy.textContent = zdanieOPracy(praca);
    wskaznikPracy.dataset.pracuje = praca === null ? 'nie' : 'tak';

    const okno = stan.okno();
    stanOkna.textContent = zdanieOStanieOkna(okno);
    stanOkna.dataset.stan = okno.rodzaj;
    element.dataset.gotowe = String(oknoGotowe(okno));

    const wToku = stan.wToku();
    wyslij.textContent = wToku ? 'Wysyłam…' : 'Wyślij';
    // Przycisk nie jest wygaszany: powtórne naciśnięcie w trakcie wysyłki
    // odpowie zdaniem w historii, a nie ciszą.
    element.dataset.wToku = String(wToku);
  }

  const przestanObserwowac = stan.obserwuj(odswiez);
  odswiez();

  return {
    element,

    pokaz() {
      element.hidden = false;
      void stan.ustalOkno();
      pole.focus();
    },

    schowaj() {
      element.hidden = true;
    },

    widoczny: () => !element.hidden,

    rozlacz() {
      przestanObserwowac();
    },
  };
}

/** Jeden wiersz rozmowy. */
function wierszWypowiedzi(wypowiedz: Wypowiedz): HTMLElement {
  const wiersz = document.createElement('article');
  wiersz.className = 'ap-wypowiedz';
  wiersz.dataset.rodzaj = wypowiedz.rodzaj;

  const kto = document.createElement('span');
  kto.className = 'ap-wypowiedz__kto';
  kto.textContent = ETYKIETY[wypowiedz.rodzaj];

  const tresc = document.createElement('p');
  tresc.className = 'ap-wypowiedz__tresc';
  tresc.textContent = wypowiedz.tresc;

  wiersz.append(kto, tresc);
  return wiersz;
}

/**
 * Zdanie o rozbieżności między profilem a tym, co dymek właśnie robi.
 *
 * Pusty napis znaczy „profil i widok mówią to samo" — wtedy pas nie stoi.
 * Napis niepusty znaczy, że ktoś zmienił profil, a dymek nie ma prawa udawać,
 * że tej zmiany nie widzi.
 */
function zdanieORozbieznosci(postac: string, ileOkien: number): string {
  if (ileOkien === 0) {
    return (
      'Profil modułu Assistant mówi dziś, że ten moduł NIE PROWADZI rozmowy ' +
      '(liczbaOkienRozmowy = 0), a mimo to dymek stoi otwarty. Zanim wpiszesz ' +
      'polecenie, sprawdź profil w okno-komunikacji/profile-wiedza.ts.'
    );
  }
  if (postac !== 'dymek-glosowy') {
    return (
      `Profil modułu Assistant mówi „postacRozmowy: ${postac}", a to jest dymek głosowy. ` +
      'Rozbieżność jest w profilu albo tutaj — nie w tym, co widzisz.'
    );
  }
  return '';
}
