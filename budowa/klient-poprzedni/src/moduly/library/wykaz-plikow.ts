import type { LibraryFile } from '../../../../shared/contract';
import type { StanBiblioteki } from './stan-biblioteki';
import { zbudujWidok } from './widoki-wykazu';

/**
 * Nawigacja po strukturze plików i zasobów wiedzy w oknie Library Explorer.
 *
 * Struktura pochodzi z danych: kontrakt nie ma komendy katalogu folderów,
 * niesie za to pole `path` każdego pliku. Wykaz katalogów składa się więc
 * z pierwszych członów ścieżek zwróconych przez rdzeń, a zbiór pusty daje
 * jedną pozycję „cały zbiór".
 *
 * Formę prezentacji rozstrzyga widok wybrany w oknie (`widoki-wykazu.ts`).
 * Wykaz nie zna żadnej z pięciu form — składa katalogi, oddaje zbiór widoczny
 * i osadza to, co widok zbudował.
 */
export interface WykazPlikow {
  element: HTMLElement;
  /**
   * Przerysowuje katalogi i ciało wykazu ze stanu modułu; oddaje zdanie
   * widoku, który nie ma czego pokazać mimo niepustego wykazu.
   */
  odswiez(): string;
}

export function utworzWykazPlikow(stan: StanBiblioteki): WykazPlikow {
  const katalogi = document.createElement('div');
  katalogi.className = 'ml-katalogi';

  const cialo = document.createElement('div');
  cialo.className = 'ml-wykaz__cialo';

  const element = document.createElement('div');
  element.className = 'ml-wykaz';
  element.append(katalogi, cialo);

  function odswiezKatalogi(): void {
    const sciezki = katalogiZbioru(stan.pliki());
    katalogi.replaceChildren(
      przyciskKatalogu('cały zbiór', '', stan),
      ...sciezki.map((sciezka) => przyciskKatalogu(sciezka, sciezka, stan)),
    );
  }

  return {
    element,
    odswiez() {
      odswiezKatalogi();
      const czynny = stan.czynny();
      const wynik = zbudujWidok(stan.widok(), stan.widoczne(), {
        zaznaczone: stan.zaznaczone(),
        czynny: czynny === null ? null : czynny.id,
        tresc: (idPliku) => stan.tresc(idPliku),
        znaczenie: (idPliku) => stan.znaczenie(idPliku),
        naZaznaczenie: (idPliku) => stan.przelaczZaznaczenie(idPliku),
        naWskazanie: (idPliku) => stan.wskaz(idPliku),
      });
      cialo.dataset['widok'] = stan.widok();
      cialo.replaceChildren(wynik.element);
      return wynik.brak;
    },
  };
}

/** Pierwsze człony ścieżek zbioru, uporządkowane po polsku. */
function katalogiZbioru(pliki: readonly LibraryFile[]): readonly string[] {
  const znane = new Set<string>();
  for (const plik of pliki) {
    const sciezka = plik.path ?? '';
    if (sciezka === '') continue;
    const czlon = sciezka.split('/')[0] ?? '';
    if (czlon !== '') znane.add(czlon);
  }
  return [...znane].sort((pierwszy, drugi) => pierwszy.localeCompare(drugi, 'pl'));
}

/** Pozycja nawigacji po katalogach; zaznaczona odpowiada zawężeniu wykazu. */
function przyciskKatalogu(etykieta: string, sciezka: string, stan: StanBiblioteki): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  const wybrany = stan.katalog() === sciezka;
  element.className = wybrany
    ? 'dn-btn dn-btn--sm dn-btn--wybrany ml-katalogi__pozycja'
    : 'dn-btn dn-btn--sm dn-btn--duch ml-katalogi__pozycja';
  element.textContent = etykieta;
  element.dataset['katalog'] = sciezka;
  element.addEventListener('click', () => stan.ustawKatalog(sciezka));
  return element;
}
