/**
 * Pytanie o jedną nazwę — modal z jednym polem tekstowym. Pobiera od Operatora
 * jeden napis albo odmowę; korzystają z niego komendy `session.rename`,
 * `session.copy` i `session.project.set`.
 *
 * Natywny `<dialog>`, nie własna nakładka: warstwę tła daje `::backdrop`
 * z `komponenty/nakladka.css`, a stos okien, pułapkę ogniska i zamknięcie
 * klawiszem Escape zapewnia przeglądarka.
 *
 * Zamknięcie oznacza odmowę, nie nazwę pustą: Escape i przycisk „Anuluj” dają
 * `null`, więc wywołujący odróżnia odmowę od świadomie pustego napisu.
 *
 * Modal jest doklejany na czas pytania i usuwany po odpowiedzi, żeby pytania
 * nie nawarstwiały się w drzewie przy każdym wierszu wykazu.
 */

export interface PytanieONazwe {
  /** Nagłówek modalu — czego dotyczy pytanie. */
  tytul: string;
  /** Etykieta pola. */
  etykieta: string;
  /** Zdanie pod polem: co się stanie po zatwierdzeniu. */
  opis?: string;
  /** Treść wstawiona do pola na start; puste pole przy pominięciu. */
  wartosc?: string;
  /** Napis przycisku zatwierdzającego. */
  napisZatwierdzenia: string;
}

/** Zwraca podany napis albo `null`, gdy Operator odmówił. */
export function zapytajONazwe(pytanie: PytanieONazwe): Promise<string | null> {
  const modal = document.createElement('dialog');
  modal.className = 'dn-modal';

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-modal-naglowek';
  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = pytanie.tytul;
  naglowek.append(tytul);

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo';
  const pole = document.createElement('label');
  pole.className = 'dn-pole';
  const etykieta = document.createElement('span');
  etykieta.className = 'dn-pole-etykieta';
  etykieta.textContent = pytanie.etykieta;
  const kontrolka = document.createElement('input');
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.type = 'text';
  kontrolka.value = pytanie.wartosc ?? '';
  pole.append(etykieta, kontrolka);
  if (pytanie.opis !== undefined) {
    const opis = document.createElement('span');
    opis.className = 'dn-pole-opis';
    opis.textContent = pytanie.opis;
    pole.append(opis);
  }
  cialo.append(pole);

  const stopka = document.createElement('div');
  stopka.className = 'dn-modal-stopka';
  const anuluj = document.createElement('button');
  anuluj.type = 'button';
  anuluj.className = 'dn-btn dn-btn--zarys';
  anuluj.textContent = 'Anuluj';
  const zatwierdz = document.createElement('button');
  zatwierdz.type = 'button';
  zatwierdz.className = 'dn-btn';
  zatwierdz.textContent = pytanie.napisZatwierdzenia;
  stopka.append(anuluj, zatwierdz);

  modal.append(naglowek, cialo, stopka);
  document.body.append(modal);
  modal.showModal();
  kontrolka.focus();
  kontrolka.select();

  return new Promise<string | null>((rozwiaz) => {
    let odpowiedz: string | null = null;

    // Jedno miejsce sprzątające: `close` przychodzi zarówno od przycisków,
    // jak i od Escape, więc usunięcie z drzewa nie ma dwóch dróg.
    modal.addEventListener('close', () => {
      modal.remove();
      rozwiaz(odpowiedz);
    });

    anuluj.addEventListener('click', () => modal.close());
    zatwierdz.addEventListener('click', () => {
      odpowiedz = kontrolka.value.trim();
      modal.close();
    });
    kontrolka.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key !== 'Enter') return;
      zdarzenie.preventDefault();
      zatwierdz.click();
    });
  });
}
