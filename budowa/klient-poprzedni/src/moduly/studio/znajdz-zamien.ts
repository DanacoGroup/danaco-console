import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { poleLogiczne, poleTekstowe, przycisk } from '../../modele/kontrolki-formularza';
import type { StanStudio } from './stan-studio';
import {
  otoczenieTrafienia,
  zamienWszystkie,
  znajdzTrafienia,
  type NastawyWyszukiwania,
} from './wyszukiwanie-tekstu';

/**
 * Znajdź/Zamień w treści bieżącego dokumentu — panel warstwy trzeciej edytora.
 *
 * Panel pracuje na buforze edytora i rdzenia nie woła. Wyszukiwanie po treści
 * WERSJI stoi w Diff/Grep Panelu i jest inną czynnością — panel mówi to
 * w objaśnieniu, żeby dwa pola wyszukiwania w jednym module nie wyglądały na
 * powielenie.
 *
 * Zamiana ma podgląd przed zatwierdzeniem, bo tego żąda opracowanie („Zamień
 * wszystko" prezentuje podgląd zmian przed zatwierdzeniem). Podgląd nie jest
 * ozdobą: zamiana wyrażeniem regularnym potrafi przepisać dokument inaczej, niż
 * Operator zamierzał, a cofnięcia po zapisie już nie ma.
 */
export interface ZnajdzZamien {
  element: HTMLElement;
  /** Przelicza trafienia po zmianie treści dokonanej poza panelem. */
  odswiez(): void;
}

const OBJASNIENIE =
  'Panel szuka w treści, którą właśnie redagujesz — w buforze edytora, także przed zapisem. ' +
  'Wyszukiwanie po treści ZAPISANYCH wersji jest osobną czynnością i stoi w Diff/Grep Panelu ' +
  '(komenda studio.diff.compare). Zamiana pokazuje podgląd i zmienia treść dopiero po ' +
  'zatwierdzeniu; wersję zakłada dopiero zapis w Studio Editorze.';

/** Górna granica wierszy podglądu — dłuższy wykaz przestaje być podglądem. */
const GRANICA_PODGLADU = 20;

export function utworzZnajdzZamien(stan: StanStudio): ZnajdzZamien {
  const wzorzec = poleTekstowe({
    etykieta: 'Znajdź w treści',
    podpowiedz: 'fraza albo wyrażenie regularne',
  });
  const zamiennik = poleTekstowe({
    etykieta: 'Zamień na',
    podpowiedz: 'treść wstawiana w miejsce trafienia',
  });
  const regularne = poleLogiczne({ etykieta: 'Wzorzec jest wyrażeniem regularnym' });
  const wielkoscLiter = poleLogiczne({ etykieta: 'Rozróżniaj wielkość liter' });

  const zamien = przycisk('Zamień wszystko', 'dn-btn dn-btn--sm dn-btn--atrament');
  zamien.dataset['czynnosc'] = 'zamien';

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dn-pole-opis ms-szukanie__podsumowanie';

  const podglad = document.createElement('ul');
  podglad.className = 'ms-szukanie__podglad';

  const element = document.createElement('div');
  element.className = 'ms-szukanie';
  element.setAttribute('aria-label', 'Znajdź i zamień w treści dokumentu');
  element.append(
    wzorzec.element,
    zamiennik.element,
    regularne.element,
    wielkoscLiter.element,
    zamien,
    utworzDymekObjasnienia(OBJASNIENIE, { powloka: 'ms-dymek', znak: 'ms-dymek__znak' }),
    podsumowanie,
    podglad,
  );

  function nastawy(): NastawyWyszukiwania {
    return {
      wzorzec: wzorzec.kontrolka.value,
      regularne: regularne.kontrolka.checked,
      wielkoscLiter: wielkoscLiter.kontrolka.checked,
    };
  }

  function odswiez(): void {
    const tresc = stan.trescRobocza();
    const nastawa = nastawy();
    if (nastawa.wzorzec === '') {
      podsumowanie.textContent =
        'Wpisz wzorzec — panel przelicza trafienia na bieżąco i pokazuje je poniżej przed zamianą.';
      podglad.replaceChildren();
      element.dataset['stan'] = 'bez-wzorca';
      return;
    }
    const wynik = znajdzTrafienia(tresc, nastawa);
    if (wynik.powod !== '') {
      // Zła składnia wzorca jest powodem nazwanym, nie pustką: bez tego zdania
      // literówka Operatora wyglądałaby jak „nic nie znaleziono".
      podsumowanie.textContent = wynik.powod;
      podglad.replaceChildren();
      element.dataset['stan'] = 'wzorzec-bledny';
      return;
    }
    if (wynik.trafienia.length === 0) {
      podsumowanie.textContent = 'Wzorca nie ma w treści dokumentu — zamiana nie miałaby czego podmienić.';
      podglad.replaceChildren();
      element.dataset['stan'] = 'bez-trafien';
      return;
    }
    const pokazane = wynik.trafienia.slice(0, GRANICA_PODGLADU);
    podsumowanie.textContent =
      wynik.trafienia.length > pokazane.length
        ? `Trafień ${wynik.trafienia.length}; poniżej pierwszych ${pokazane.length}. Zamiana obejmie wszystkie.`
        : `Trafień ${wynik.trafienia.length}. Zamiana obejmie wszystkie.`;
    podglad.replaceChildren(
      ...pokazane.map((trafienie) => {
        const polozenie = document.createElement('span');
        polozenie.className = 'dn-plakietka ms-szukanie__polozenie';
        polozenie.textContent = `znak ${trafienie.poczatek}`;

        const otoczenie = document.createElement('code');
        otoczenie.className = 'ms-szukanie__otoczenie';
        otoczenie.textContent = otoczenieTrafienia(tresc, trafienie);

        const pozycja = document.createElement('li');
        pozycja.append(polozenie, otoczenie);
        return pozycja;
      }),
    );
    element.dataset['stan'] = 'trafienia';
  }

  for (const kontrolka of [wzorzec.kontrolka, regularne.kontrolka, wielkoscLiter.kontrolka]) {
    kontrolka.addEventListener('input', odswiez);
  }

  zamien.addEventListener('click', () => {
    const wynik = zamienWszystkie(stan.trescRobocza(), nastawy(), zamiennik.kontrolka.value);
    if (wynik.powod !== '') {
      podsumowanie.textContent = wynik.powod;
      return;
    }
    if (wynik.liczba === 0) {
      podsumowanie.textContent = 'Nie było czego zamienić — wzorca nie ma w treści dokumentu.';
      return;
    }
    // Treść idzie przez stan, nie wprost do kontrolki: jeden dokument na cały
    // moduł, więc pozostałe okna mają się dowiedzieć.
    stan.ustawTresc(wynik.tresc);
    podsumowanie.textContent =
      `Zamieniono trafień: ${wynik.liczba}. Zmiana stoi w buforze edytora — zapis w Studio ` +
      'Editorze utrwali ją i założy wersję.';
  });

  odswiez();
  return { element, odswiez };
}
