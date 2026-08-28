import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';
import { odczytajBase64 } from './tresc-base64';

/**
 * Dołożenie wersji dokumentu w Versioning Panelu jest jedynym wejściem, którym
 * rośnie historia dokumentu: wersja niesie albo nową treść, albo etykietę
 * kamienia milowego postawionego na treści bieżącej.
 */
export interface DolozenieWersji {
  element: HTMLElement;
}

/** Sprawca zmiany wykonanej ręką człowieka przy oknie modułu, odróżniający wpis od wersji wytworzonych przez model. */
const SPRAWCA = 'Operator';

export function utworzDolozenieWersji(
  stan: StanBiblioteki,
  poZmianie: () => void,
): DolozenieWersji {
  const odpowiedz = utworzWierszOdpowiedzi();

  const etykieta = poleTekstowe({
    etykieta: 'Etykieta wersji',
    podpowiedz: 'np. wydanie 1.0',
    opis: 'Bez etykiety rdzeń nazwie wersję jej kolejnością w historii dokumentu.',
  });

  const wybor = document.createElement('input');
  wybor.type = 'file';
  wybor.hidden = true;
  wybor.addEventListener('change', () => {
    const plik = wybor.files?.[0];
    if (plik === undefined) return;
    void dolozZTrescia(plik);
  });

  const zPliku = przycisk('Dołóż wersję z pliku', 'dn-btn dn-btn--sm dn-btn--atrament');
  zPliku.dataset['czynnosc'] = 'wersja-dolozenie';
  zPliku.addEventListener('click', () => wybor.click());

  const kamien = przycisk('Oznacz kamień milowy', 'dn-btn dn-btn--sm dn-btn--zarys');
  kamien.dataset['czynnosc'] = 'wersja-kamien';
  kamien.addEventListener('click', () => void dolozZnacznik());

  async function dolozZTrescia(plik: File): Promise<void> {
    const dokument = stan.czynny();
    if (dokument === null) {
      odpowiedz.pokaz('Wskaż dokument w Library Explorer — wersję dokłada się do pliku.', false);
      return;
    }
    odpowiedz.pokaz(`Odczyt pliku „${plik.name}" przed wysłaniem…`, true);
    const tresc = await odczytajBase64(plik);
    if (tresc === null) {
      odpowiedz.pokaz(`Odczyt pliku „${plik.name}" nie powiódł się — nic nie zostało wysłane.`, false);
      return;
    }
    await wyslij({ fileId: dokument.id, contentBase64: tresc }, `Wersja z pliku „${plik.name}"`);
  }

  async function dolozZnacznik(): Promise<void> {
    const dokument = stan.czynny();
    if (dokument === null) {
      odpowiedz.pokaz('Wskaż dokument w Library Explorer — kamień milowy stawia się na pliku.', false);
      return;
    }
    await wyslij({ fileId: dokument.id }, 'Kamień milowy');
  }

  /** Jedna droga wysyłki dla obu przycisków — jeden opis powodzenia i odmowy. */
  async function wyslij(zadanie: { fileId: string; contentBase64?: string }, czynnosc: string): Promise<void> {
    const nazwa = etykieta.kontrolka.value.trim();
    odpowiedz.pokaz(`${czynnosc} — rdzeń dokłada wersję…`, true);
    const wynik = await stan.zrodlo.dolozWersje({
      ...zadanie,
      ...(nazwa === '' ? {} : { label: nazwa }),
      author: SPRAWCA,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy(czynnosc, wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const wersja = wynik.wynik.version;
    const plik = wynik.wynik.file;
    const nazwaWersji = wersja.label ?? wersja.id;
    stan.wchlon(plik);
    etykieta.kontrolka.value = '';
    poZmianie();

    // Wskazanie wersji niesionej przez dokument stoi w polu pliku, a nie w domyśle okna.
    if (plik.versionId !== wersja.id) {
      odpowiedz.pokaz(
        `${czynnosc}: rdzeń założył wersję ${nazwaWersji} (${wersja.id}), ale dokument ` +
          `oddał ze wskazaniem wersji ${plik.versionId ?? 'nieokreślonej'}. Nie potwierdzam ` +
          'dołożenia — zamówienie i odpowiedź mówią o czym innym.',
        false,
      );
      return;
    }
    odpowiedz.pokaz(`${czynnosc}: wpis ${nazwaWersji} założony — pytam rdzeń o treść…`, true);
    const stanTresci = await stan.zbadajTresc(plik.id);
    if (stanTresci.werdykt === 'brak') {
      odpowiedz.pokaz(
        `${czynnosc}: wersja ${nazwaWersji} stoi w historii, ale dokument nadal nie ma ` +
          `treści w repozytorium — ${stanTresci.powod}. „Przywróć" nie ma dokąd wrócić.`,
        false,
      );
      return;
    }
    if (stanTresci.werdykt === 'odwolanie') {
      odpowiedz.pokaz(
        `${czynnosc}: wersja ${nazwaWersji} stoi w historii i dokument ją niesie, ale ` +
          `podgląd oddał wskazanie zamiast treści — ${stanTresci.powod}.`,
        false,
      );
      return;
    }
    // Treść dostępna wolno napisać wyłącznie po stwierdzeniu osiągalności, nie po odmowie odczytu.
    if (stanTresci.werdykt !== 'osiagalna') {
      odpowiedz.pokaz(
        `${czynnosc}: wersja ${nazwaWersji} stoi w historii, ale rdzeń nie odpowiedział ` +
          `o treści dokumentu — ${stanTresci.powod}. Czy „Przywróć" ma dokąd wrócić, ` +
          'nie wiadomo.',
        false,
      );
      return;
    }
    // Sumy kontrolne z jednej odpowiedzi mówią, czy treść oddana przez podgląd należy do tej wersji.
    if (wersja.checksum !== undefined && wersja.checksum !== plik.checksum) {
      odpowiedz.pokaz(
        `${czynnosc}: dołożona jako ${nazwaWersji}, rdzeń oddaje treść dokumentu — ale suma ` +
          `kontrolna wersji (${wersja.checksum}) różni się od sumy dokumentu ` +
          `(${plik.checksum ?? 'brak'}), więc nie jest to treść tej wersji.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      `${czynnosc}: dołożona jako ${nazwaWersji}; dokument niesie tę wersję, a rdzeń oddaje ` +
        'jej treść do podglądu.',
      true,
    );
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-wersje__pasek';
  pasek.append(zPliku, kamien, wybor);

  const element = document.createElement('div');
  element.className = 'ml-dolozenie';
  element.append(etykieta.element, pasek, odpowiedz.element);

  return { element };
}
