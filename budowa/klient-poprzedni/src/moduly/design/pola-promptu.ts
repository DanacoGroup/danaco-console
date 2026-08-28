import { DesignAssetKind, type Channel, type DesignAsset, type DesignPrompt } from '../../../../shared/contract';
import { poleTekstowe, poleWyboru, ustawPozycje } from '../../modele/kontrolki-formularza';
import { dopnijDymek } from './dymek-objasnienia';
import { czyKanalObrazowy, powodNieprzydatnosci } from './kanaly-obrazowe';
import { nazwaRodzaju, nazwaZasobu } from './karta-zasobu';
import { przytnijKod } from './przyciecie-pol';
import { utworzSuwaki, type SuwakiPromptu } from './suwaki-promptu';

/**
 * Pola strukturalne Prompt Buildera — siedem pól promptu, silnik i referencja, każde odwzorowujące
 * jedno pole promptu kontraktu.
 */
export interface PolaPromptu {
  element: HTMLElement;
  /** Prompt złożony z pól — gotowa treść żądania `design.asset.generate`. */
  prompt(): DesignPrompt;
  /** Zasób referencyjny generowania obraz-do-obrazu; pusty znaczy brak. */
  idReferencji(): string;
  /** Kanał obrazowy wskazany przez Operatora; pusty znaczy, że Operator nie wskazał, nie że kanał pusty. */
  idKanalu(): string;
  /** Rodzaj zasobu wskazany przez Operatora — pole `kind` żądania; pusty znaczy brak wskazania. */
  rodzaj(): DesignAssetKind | '';
  /** Nanosi prompt na pola — droga powrotu z historii poleceń. */
  nanies(prompt: DesignPrompt): void;
  /** Odświeża wykazy silników, zasobów referencyjnych i podpowiedzi stylów. */
  odswiez(silniki: readonly Channel[], zasoby: readonly DesignAsset[], style: readonly string[]): void;
  /** Liczba znaków treści promptu — licznik długości z panelu akcji. */
  dlugosc(): number;
}

/**
 * Nazwa stera silnika bez wartości bieżącej — doklejana przy każdym odczycie.
 *
 * Stoi osobno, bo etykieta niesie nastawę biegnącą, a nie samą nazwę pola.
 * Nazwa i nastawa pochodzą z jednego miejsca, inaczej rozjadą się przy
 * pierwszej zmianie.
 */
export const ETYKIETA_SILNIKA = 'Silnik generujący (kanał obrazowy)';

/** Nazwa stera rodzaju zasobu bez wartości bieżącej, doklejana przy każdym odczycie etykiety sterownika. */
export const ETYKIETA_RODZAJU = 'Rodzaj zasobu';

/** Rodzaje zasobu wybieralne w kreatorze — te, które generowanie potrafi wytworzyć z kanału obrazowego. */
const RODZAJE_DO_WYBORU: readonly DesignAssetKind[] = [
  DesignAssetKind.Image,
  DesignAssetKind.Vector,
  DesignAssetKind.Document,
];

/** Pole strukturalne: klucz pola kontraktu, etykieta widoczna, podpowiedź w kontrolce, objaśnienie w dymku. */
const POLA: readonly (readonly [keyof DesignPrompt, string, string, string])[] = [
  ['subject', 'Temat', 'co ma przedstawiać zasób', 'Pole subject — jedyne wymagane pole promptu.'],
  ['style', 'Styl', 'np. rysunek techniczny', 'Pole style. Podpowiedź składa się ze stylów użytych w tej sesji.'],
  ['composition', 'Kompozycja', 'np. kadr centralny', 'Pole composition — rozmieszczenie elementów w kadrze.'],
  ['lighting', 'Oświetlenie', 'np. światło boczne', 'Pole lighting — charakter światła sceny.'],
  ['palette', 'Paleta barw', 'np. atrament i sygnał', 'Pole palette — zakres barw zasobu.'],
  ['aspectRatio', 'Proporcje kadru', 'np. 16:9', 'Pole aspectRatio — stosunek boków zasobu.'],
  ['exclusions', 'Wykluczenia', 'czego ma nie być', 'Pole exclusions — czego generowanie ma uniknąć.'],
];

export function utworzPolaPromptu(): PolaPromptu {
  const kontrolki = new Map<keyof DesignPrompt, HTMLInputElement>();
  const podpowiedziStylow = document.createElement('datalist');
  podpowiedziStylow.id = 'md-style-podpowiedzi';

  const element = document.createElement('div');
  element.className = 'md-prompt__pola';

  for (const [klucz, etykieta, podpowiedz, objasnienie] of POLA) {
    const pole = poleTekstowe({ etykieta, podpowiedz });
    dopnijDymek(pole.element, objasnienie);
    if (klucz === 'style') pole.kontrolka.setAttribute('list', podpowiedziStylow.id);
    kontrolki.set(klucz, pole.kontrolka);
    element.append(pole.element);
  }

  // Wskazanie silnika rozstrzyga, którym kanałem powstanie zasób; rdzeń odmawia kanału tekstowego.
  const silnik = poleWyboru({
    etykieta: ETYKIETA_SILNIKA,
    opis:
      'Kanały obrazowe z rejestru kanałów modelu (channel.list, adapter „obrazy"). Wskazanie ' +
      'idzie do rdzenia polem channelId żądania design.asset.generate i ROZSTRZYGA, którym ' +
      'kanałem powstanie zasób. Bez wskazania rdzeń bierze pierwszy czynny kanał obrazowy konta.',
  }, []);
  dopnijDymek(
    silnik.element,
    'Pole channelId żądania. Kanał tekstowy jest odmawiany przez rdzeń przed wysłaniem ' +
      'polecenia — fragment tekstu nie jest obrazem — więc wykaz podaje kanały nieobrazowe ' +
      'jako niewybieralne wraz z powodem, zamiast je ukrywać.',
  );
  const etykietaSilnika = silnik.element.querySelector('.dn-pole-etykieta');

  // Rdzeń bez wskazanego rodzaju zakłada zasób obrazu; wartość pusta zostawia rozstrzygnięcie rdzeniowi.
  const rodzaj = poleWyboru({
    etykieta: ETYKIETA_RODZAJU,
    opis:
      'Pole kind żądania design.asset.generate. Bez wskazania rdzeń zakłada zasób rodzaju ' +
      'obraz rastrowy. Rodzaj opisuje ZASÓB, nie kanał — kanał obrazowy oddaje bajty, ' +
      'a rdzeń mierzy ich format z nagłówka utrwalonego pliku.',
  }, [
    { wartosc: '', etykieta: 'rodzaj rozstrzyga rdzeń (obraz rastrowy)' },
    ...RODZAJE_DO_WYBORU.map((wartosc) => ({ wartosc, etykieta: nazwaRodzaju(wartosc) })),
  ]);
  dopnijDymek(rodzaj.element, 'Wartość pola kind. Pominięcie pola jest w kontrakcie opisane jako rodzaj obrazu.');
  const etykietaRodzaju = rodzaj.element.querySelector('.dn-pole-etykieta');

  const referencja = poleWyboru({
    etykieta: 'Obraz referencyjny (obraz-do-obrazu)',
    opis: 'Zasób z Assets Panel; pole referenceAssetId żądania design.asset.generate.',
  }, []);
  dopnijDymek(referencja.element, 'Wskazany zasób staje się punktem wyjścia generowania zamiast czystego promptu.');

  const suwaki: SuwakiPromptu = utworzSuwaki();

  element.append(podpowiedziStylow, silnik.element, rodzaj.element, referencja.element, suwaki.element);

  /** Przepisuje nastawę bieżącą na etykiety obu sterów; etykieta czyta wykaz przy każdej zmianie. */
  function przepiszNastawy(): void {
    if (etykietaSilnika !== null) {
      etykietaSilnika.textContent = `${ETYKIETA_SILNIKA}: ${nastawaWyboru(silnik.kontrolka)}`;
    }
    if (etykietaRodzaju !== null) {
      etykietaRodzaju.textContent = `${ETYKIETA_RODZAJU}: ${nastawaWyboru(rodzaj.kontrolka)}`;
    }
  }

  silnik.kontrolka.addEventListener('change', przepiszNastawy);
  rodzaj.kontrolka.addEventListener('change', przepiszNastawy);
  przepiszNastawy();

  // Pole promptu przycinane jest tak samo jak etykieta — wspólną notacją dla przeglądarki i rdzenia.
  function tekst(klucz: keyof DesignPrompt): string {
    return przytnijKod(kontrolki.get(klucz)?.value ?? '');
  }

  return {
    element,

    prompt() {
      const prompt: DesignPrompt = { subject: tekst('subject'), ...suwaki.wartosci() };
      for (const [klucz] of POLA) {
        if (klucz === 'subject') continue;
        const wartosc = tekst(klucz);
        if (wartosc !== '') Object.assign(prompt, { [klucz]: wartosc });
      }
      if (silnik.kontrolka.value !== '') prompt.engine = silnik.kontrolka.value;
      return prompt;
    },

    idReferencji: () => referencja.kontrolka.value,

    idKanalu: () => silnik.kontrolka.value,

    // Wartość pola pochodzi z wykazu rodzajów, więc jest jedną z nich albo pustką, nigdy inną treścią.
    rodzaj: () => rodzaj.kontrolka.value as DesignAssetKind | '',

    nanies(prompt) {
      for (const [klucz] of POLA) {
        const kontrolka = kontrolki.get(klucz);
        if (kontrolka === undefined) continue;
        const wartosc = prompt[klucz];
        kontrolka.value = typeof wartosc === 'string' ? wartosc : '';
      }
      // Kanał z historii wpisywany jest wprost; etykieta czyta ster, więc nie rozejdzie się z nastawą.
      if (prompt.engine !== undefined) silnik.kontrolka.value = prompt.engine;
      suwaki.nanies(prompt);
      przepiszNastawy();
    },

    odswiez(silniki, zasoby, style) {
      // Wykaz zawężony do kanałów obrazowych, ale nie oczyszczony: tekstowy zostaje widoczny, niewybieralny.
      const obrazowe = silniki.filter(czyKanalObrazowy);
      const pozostale = silniki.filter((kanal) => !czyKanalObrazowy(kanal));
      ustawPozycje(silnik.kontrolka, [
        { wartosc: '', etykieta: zdanieWyjsciowegoWyboru(obrazowe.length) },
        ...obrazowe.map((kanal) => ({ wartosc: kanal.id, etykieta: `${kanal.name} · obrazowy` })),
      ]);
      for (const kanal of pozostale) {
        const pozycja = document.createElement('option');
        pozycja.value = kanal.id;
        pozycja.disabled = true;
        pozycja.textContent = `${kanal.name} — ${powodNieprzydatnosci(kanal)}`;
        silnik.kontrolka.append(pozycja);
      }
      ustawPozycje(referencja.kontrolka, [
        { wartosc: '', etykieta: 'bez obrazu referencyjnego' },
        ...zasoby.map((zasob) => ({ wartosc: zasob.id, etykieta: nazwaZasobu(zasob) })),
      ]);
      podpowiedziStylow.replaceChildren(...style.map(pozycjaPodpowiedzi));
      przepiszNastawy();
    },

    dlugosc() {
      return POLA.reduce((suma, [klucz]) => suma + tekst(klucz).length, 0);
    },
  };
}

function pozycjaPodpowiedzi(styl: string): HTMLOptionElement {
  const element = document.createElement('option');
  element.value = styl;
  return element;
}

/**
 * Nastawa bieżąca stera — treść zaznaczonej pozycji, nigdy jej wartość surowa.
 *
 * Operator czyta etykietę, a nie identyfikator kanału; wpisanie tu `kanal.id`
 * dałoby etykietę mówiącą co innego niż rozwinięty wykaz.
 */
function nastawaWyboru(kontrolka: HTMLSelectElement): string {
  // Pozycja brana przez indeks zaznaczenia, nie przez zbiór opcji, bywający jeszcze sprzed zmiany.
  const zaznaczona = kontrolka.options[kontrolka.selectedIndex];
  if (zaznaczona === undefined) return 'wykaz pusty';
  return zaznaczona.textContent ?? '';
}

/**
 * Pozycja wyjściowa wykazu silników — mówi, czy jest z czego wybierać.
 *
 * Rejestr bez ani jednego kanału obrazowego jest stanem, po którym generowanie
 * odmówi — Operator widzi to przed naciśnięciem „Generuj".
 */
function zdanieWyjsciowegoWyboru(obrazowych: number): string {
  if (obrazowych === 0) {
    return 'brak kanału obrazowego w rejestrze — generowanie odmówi';
  }
  return `kanał rozstrzyga rdzeń (pierwszy czynny z ${obrazowych})`;
}
