/**
 * Rusztowanie podglądowe warstwy rozmowy, uruchamiane poza pakietem produktu i pokazujące
 * ją bez powłoki ani nawigacji modułów.
 */

import '../motyw/motyw.css';
import '../komponenty/indeks.css';

import { uruchomPodglad } from './polaczenie-podgladu';

const scena = document.getElementById('scena');
const stan = document.getElementById('stan');
const przelacznik = document.getElementById('motyw');

if (scena instanceof HTMLElement && stan instanceof HTMLElement) {
  uruchomPodglad({ scena, stan });
}

// Oba motywy interfejsu są równoprawne w tym podglądzie — kliknięcie przełącznika pokazuje
// kolejno jeden i drugi.
przelacznik?.addEventListener('click', () => {
  const biezacy = document.documentElement.dataset['theme'] === 'dark' ? 'light' : 'dark';
  document.documentElement.dataset['theme'] = biezacy;
  przelacznik.textContent = biezacy === 'dark' ? 'Motyw jasny' : 'Motyw ciemny';
});
