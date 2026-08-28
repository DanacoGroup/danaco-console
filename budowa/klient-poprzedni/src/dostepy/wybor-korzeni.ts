/**
 * Wybór korzeni katalogowych nadania — podzbiór korzeni punktu dostępu. Pusty zbiór
 * korzeni nadania znaczy w kontrakcie komplet korzeni punktu, a nie brak dostępu,
 * i ta reguła rządzi całą kontrolką.
 */
export interface WyborKorzeni {
  /** Element osadzany w wierszu nadania albo w karcie punktu. */
  element: HTMLElement;
  /** Zaznaczone korzenie; komplet zaznaczony daje tablicę pustą. */
  odczytaj(): string[];
  /** Nanosi zaznaczenie; tablica pusta zaznacza komplet. */
  ustaw(wybrane: readonly string[]): void;
  /** Zgłasza zmianę zaznaczenia. */
  naZmiane(sluchacz: () => void): void;
}

export interface ZaleznosciWyboru {
  /** Korzenie punktu dostępu — granica, poza którą nadanie nie sięga. */
  korzenie: readonly string[];
  /** Przedrostek identyfikatorów pól, żeby dwa wykazy się nie zlały. */
  identyfikator: string;
}

export function utworzWyborKorzeni(zaleznosci: ZaleznosciWyboru): WyborKorzeni {
  const { korzenie, identyfikator } = zaleznosci;

  const element = document.createElement('fieldset');
  element.className = 'dd-korzenie';

  const legenda = document.createElement('legend');
  legenda.className = 'dd-korzenie__legenda';
  legenda.textContent = 'Korzenie katalogowe';

  const sluchacze: Array<() => void> = [];
  const pola: HTMLInputElement[] = [];

  const wykaz = document.createElement('div');
  wykaz.className = 'dd-korzenie__wykaz';

  korzenie.forEach((korzen, numer) => {
    const pole = document.createElement('input');
    pole.type = 'checkbox';
    pole.className = 'dn-check dd-korzenie__pole';
    pole.id = `${identyfikator}-korzen-${numer}`;
    pole.value = korzen;
    pole.checked = true;
    pole.addEventListener('change', () => {
      odswiezPodsumowanie();
      for (const sluchacz of [...sluchacze]) sluchacz();
    });

    const etykieta = document.createElement('label');
    etykieta.className = 'dd-korzenie__etykieta';
    etykieta.htmlFor = pole.id;
    etykieta.textContent = korzen;
    etykieta.title = korzen;

    const wiersz = document.createElement('div');
    wiersz.className = 'dd-korzenie__wiersz';
    wiersz.append(pole, etykieta);

    pola.push(pole);
    wykaz.append(wiersz);
  });

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dd-korzenie__podsumowanie';

  element.append(legenda, korzenie.length > 0 ? wykaz : brakKorzeni(), podsumowanie);

  /** Zdanie pod wykazem mówi, co naprawdę obejmie nadanie. */
  function odswiezPodsumowanie(): void {
    const zaznaczone = pola.filter((pole) => pole.checked).length;
    if (korzenie.length === 0) {
      podsumowanie.textContent =
        'Punkt nie zgłasza korzeni — nadanie obejmie obszar ustalony przez sam punkt.';
      return;
    }
    if (zaznaczone === korzenie.length) {
      podsumowanie.textContent = `Nadanie obejmie komplet korzeni punktu (${korzenie.length}).`;
      return;
    }
    if (zaznaczone === 0) {
      podsumowanie.textContent =
        'Nie zaznaczono żadnego korzenia — nadanie obejmie komplet korzeni punktu.';
      return;
    }
    podsumowanie.textContent = `Nadanie obejmie ${zaznaczone} z ${korzenie.length} korzeni punktu.`;
  }

  odswiezPodsumowanie();

  return {
    element,

    odczytaj() {
      const zaznaczone = pola.filter((pole) => pole.checked).map((pole) => pole.value);
      return zaznaczone.length === pola.length ? [] : zaznaczone;
    },

    ustaw(wybrane) {
      const komplet = wybrane.length === 0;
      for (const pole of pola) {
        pole.checked = komplet || wybrane.includes(pole.value);
      }
      odswiezPodsumowanie();
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}

/**
 * Punkt bez korzeni nie daje pustego prostokąta, tylko zdanie o tym, że korzeni
 * nie wymienia. Nadanie obejmie wtedy obszar, który punkt sam potwierdza przy
 * sprawdzeniu, więc milczenie kontrolki byłoby wprowadzeniem w błąd.
 */
function brakKorzeni(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dd-korzenie__pusty';
  element.textContent =
    'Punkt dostępu nie wymienia korzeni. Nadanie obejmie obszar, który punkt sam potwierdza przy sprawdzeniu.';
  return element;
}
