import '../motyw/motyw.css';
import '../komponenty/indeks.css';

import { uruchomPodglad } from './polaczenie-podgladu';

// Rusztowanie podglądowe — poza pakietem produktu.
//
// Pokazuje warstwę rozmowy samą: bez powłoki, bez nawigacji modułów, bez pasa
// kart sesji — na gołej scenie i z przełącznikiem motywu. Dzięki temu usterkę
// rozmowy widać bez zgadywania, czy winna jest rozmowa, czy warstwa, która ją
// osadza. Osobno wykonuje `zapewnijKanalGlowny`, czyli założenie wiersza
// rejestru na świeżej bazie — drogę, którą produkt przechodzi tylko ręką
// Operatora w panelu sterowania.
//
// Podgląd nie dowodzi, że rozmowa działa w produkcie: produkt osadza ją
// w module, w powłoce i w oknie sesji, a tutaj powłoki w ogóle nie ma.
//
// Uruchomienie:  cd budowa/client && npm run dev
//                → http://localhost:5173/src/rozmowa/podglad-rozmowy.html
//                (opcjonalnie `?pytanie=…`, żeby zobaczyć pełną turę)
//
// Sam plik to wyłącznie kompozycja: arkusze, odszukanie miejsc w dokumencie,
// przełącznik motywu, uruchomienie. Logika mieszka w `polaczenie-podgladu.ts`.

const scena = document.getElementById('scena');
const stan = document.getElementById('stan');
const przelacznik = document.getElementById('motyw');

if (scena instanceof HTMLElement && stan instanceof HTMLElement) {
  uruchomPodglad({ scena, stan });
}

// Oba motywy są równoprawne — podgląd pokazuje jeden i drugi.
przelacznik?.addEventListener('click', () => {
  const biezacy = document.documentElement.dataset['theme'] === 'dark' ? 'light' : 'dark';
  document.documentElement.dataset['theme'] = biezacy;
  przelacznik.textContent = biezacy === 'dark' ? 'Motyw jasny' : 'Motyw ciemny';
});
