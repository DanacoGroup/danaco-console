// Punkt wejścia strony podglądu strony głównej.
//
// Wzorzec ten sam co w `motyw/podglad-zetonow.ts`: widok daje się obejrzeć bez
// montażu w powłoce aplikacji. Podgląd nie jest częścią produktu — nie wchodzi
// do pakietu `main.ts` i niczego z niego nie importuje.
//
// Podgląd pokazuje oba motywy, bo są równoprawne.

import '../motyw/motyw.css';
import '../komponenty/indeks.css';

import { elementIkony } from '../ikony/ikony';
import { motywObowiazujacy, przelaczMotyw, uruchomMotyw, ustawMotyw } from '../motyw/motyw';
import { POZYCJE_SRODOWISK } from './pozycje-srodowisk';
import { utworzStroneGlowna } from './indeks';

// Parametry podglądu w adresie: `#motyw=light&czynne=talkin`.
// Bez nich obowiązuje preferencja systemu i brak środowiska czynnego.
// Dzięki nim oba motywy i stan czynny karty dają się obejrzeć bez klikania —
// także z narzędzia zrzucającego obraz strony.
const parametry = new URLSearchParams(location.hash.slice(1));

uruchomMotyw();

const zadanyMotyw = parametry.get('motyw');
if (zadanyMotyw === 'light' || zadanyMotyw === 'dark') {
  ustawMotyw(zadanyMotyw);
}

const strona = utworzStroneGlowna();

const zadaneCzynne = POZYCJE_SRODOWISK.find((p) => p.kod === parametry.get('czynne'));
if (zadaneCzynne !== undefined) {
  strona.ustawSrodowiskoCzynne(zadaneCzynne.kod);
}

const dziennik = document.createElement('output');
dziennik.className = 'dn-plakietka dn-plakietka--informacja podglad__dziennik';
dziennik.textContent = 'Wybór zgłoszony przez stronę pojawi się tutaj.';

strona.naWyborSrodowiska((pozycja) => {
  strona.ustawSrodowiskoCzynne(pozycja.kod);
  dziennik.textContent = `Środowisko: ${pozycja.nazwa} (${pozycja.kod})`;
});

strona.naWyborKomponentu((pozycja) => {
  dziennik.textContent = `Komponent: ${pozycja.nazwa} — ${pozycja.wezwanie}`;
});

strona.naWyborUstawienia((pozycja) => {
  dziennik.textContent = `Ustawienia: ${pozycja.nazwa}`;
});

const przelacznik = document.createElement('button');
przelacznik.type = 'button';
przelacznik.className = 'dn-btn dn-btn--zarys dn-btn--sm';
odswiezPrzelacznik();
przelacznik.addEventListener('click', () => {
  przelaczMotyw();
  odswiezPrzelacznik();
});

/** Przełącznik nazywa motyw, do którego prowadzi, i niesie jego ikonę. */
function odswiezPrzelacznik(): void {
  const jasny = motywObowiazujacy() === 'light';
  przelacznik.replaceChildren(
    elementIkony(jasny ? 'ksiezyc' : 'slonce', { rozmiar: 16 }),
    document.createTextNode(jasny ? 'Motyw ciemny' : 'Motyw jasny'),
  );
}

const listwaPodgladu = document.createElement('div');
listwaPodgladu.className = 'podglad__listwa';
listwaPodgladu.append(przelacznik, dziennik);

document.body.append(listwaPodgladu, strona.element);
