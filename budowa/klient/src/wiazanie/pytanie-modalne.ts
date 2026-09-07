// Pytanie o wartości w oknie modalnym — sposób, którym okna platformowe pytają
// Operatora, tak jak modal zmiany hasła stojący w prototypie Ustawień.
import type { PoleSzuflady } from './czynnosci-okna.ts';

export interface PytanieModalne {
  tytul: string;
  opis?: string;
  pola: PoleSzuflady[];
  wykonanie: string;
  /** Zdanie o następstwie nieodwracalnym; wymusza drugie naciśnięcie. */
  nieodwracalne?: string;
}

/*
zapytajWOknieModalnym stawia modal z polami i czeka na wartości od Operatora.

Okna platformowe nie mają paneli z belką, więc czynność nie ma gdzie rozwinąć
szuflady — pyta modalem, tak jak robi to prototyp Ustawień.
*/
export async function zapytajWOknieModalnym(
  pytanie: PytanieModalne,
): Promise<Record<string, string> | null> {
  const dokument = document;
  dokument.querySelector('[data-pytanie-modalne]')?.remove();

  const tlo = dokument.createElement('div');
  tlo.className = 'dn-modal-tlo';
  tlo.dataset.pytanieModalne = '';
  tlo.dataset.otwarte = 'tak';
  const modal = dokument.createElement('div');
  modal.className = 'dn-modal';
  modal.setAttribute('role', 'dialog');
  modal.setAttribute('aria-modal', 'true');

  const naglowek = dokument.createElement('header');
  naglowek.className = 'dn-modal-naglowek';
  const tytul = dokument.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = pytanie.tytul;
  const zamknij = dokument.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn-ikona dn-modal-zamknij';
  zamknij.setAttribute('aria-label', 'Zamknij bez wykonania');
  zamknij.textContent = '×';
  naglowek.append(tytul, zamknij);

  const cialo = dokument.createElement('div');
  cialo.className = 'dn-modal-cialo';
  if (pytanie.opis !== undefined) {
    const zdanie = dokument.createElement('p');
    zdanie.className = 'dn-meta';
    zdanie.textContent = pytanie.opis;
    cialo.appendChild(zdanie);
  }
  const kontrolki = new Map<string, HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>();
  for (const pole of pytanie.pola) {
    const opakowanie = dokument.createElement('div');
    opakowanie.className = 'dn-pole';
    const nazwa = dokument.createElement('span');
    nazwa.className = 'dn-pole-etykieta';
    nazwa.textContent = pole.etykieta;
    const kontrolka = zbudujKontrolke(dokument, pole);
    opakowanie.append(nazwa, kontrolka);
    cialo.appendChild(opakowanie);
    kontrolki.set(pole.klucz, kontrolka);
  }

  const stopka = dokument.createElement('footer');
  stopka.className = 'dn-modal-stopka';
  const uwaga = dokument.createElement('p');
  uwaga.className = 'dn-meta';
  uwaga.hidden = true;
  const wykonanie = dokument.createElement('button');
  wykonanie.type = 'button';
  wykonanie.className = 'dn-btn dn-btn--sygnal';
  wykonanie.textContent = pytanie.wykonanie;
  stopka.append(uwaga, wykonanie);

  modal.append(naglowek, cialo, stopka);
  tlo.appendChild(modal);
  dokument.body.appendChild(tlo);
  kontrolki.values().next().value?.focus();

  return new Promise((rozwiaz) => {
    const domknij = (wynik: Record<string, string> | null): void => {
      tlo.remove();
      rozwiaz(wynik);
    };
    zamknij.addEventListener('click', () => {
      domknij(null);
    });
    tlo.addEventListener('click', (zdarzenie) => {
      if (zdarzenie.target === tlo) domknij(null);
    });
    tlo.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key === 'Escape') domknij(null);
    });
    wykonanie.addEventListener('click', () => {
      const brakujace = pytanie.pola
        .filter((pole) => pole.wymagane === true && (kontrolki.get(pole.klucz)?.value ?? '') === '')
        .map((pole) => pole.etykieta);
      for (const pole of pytanie.pola) {
        kontrolki.get(pole.klucz)?.setAttribute('aria-invalid',
          String(brakujace.includes(pole.etykieta)));
      }
      if (brakujace.length > 0) {
        uwaga.hidden = false;
        uwaga.textContent = `Bez wartości nie da się wykonać: ${brakujace.join(', ')}.`;
        return;
      }
      if (pytanie.nieodwracalne !== undefined && uwaga.dataset.odslonione !== 'tak') {
        uwaga.dataset.odslonione = 'tak';
        uwaga.hidden = false;
        uwaga.textContent = pytanie.nieodwracalne;
        wykonanie.textContent = `${pytanie.wykonanie} — potwierdź`;
        return;
      }
      const wartosci: Record<string, string> = {};
      for (const [klucz, kontrolka] of kontrolki) wartosci[klucz] = kontrolka.value.trim();
      domknij(wartosci);
    });
  });
}

function zbudujKontrolke(
  dokument: Document,
  pole: PoleSzuflady,
): HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement {
  if (pole.wybor !== undefined) {
    const wybor = dokument.createElement('select');
    wybor.className = 'dn-pole-kontrolka';
    for (const [wartosc, nazwa] of pole.wybor) {
      const pozycja = dokument.createElement('option');
      pozycja.value = wartosc;
      pozycja.textContent = nazwa;
      wybor.appendChild(pozycja);
    }
    wybor.value = pole.wartosc ?? pole.wybor[0]?.[0] ?? '';
    return wybor;
  }
  if (pole.obszerne === true) {
    const obszar = dokument.createElement('textarea');
    obszar.className = 'dn-pole-kontrolka';
    obszar.rows = 4;
    obszar.value = pole.wartosc ?? '';
    if (pole.podpowiedz !== undefined) obszar.placeholder = pole.podpowiedz;
    return obszar;
  }
  const wpis = dokument.createElement('input');
  wpis.type = 'text';
  wpis.className = 'dn-pole-kontrolka';
  wpis.value = pole.wartosc ?? '';
  if (pole.podpowiedz !== undefined) wpis.placeholder = pole.podpowiedz;
  return wpis;
}
