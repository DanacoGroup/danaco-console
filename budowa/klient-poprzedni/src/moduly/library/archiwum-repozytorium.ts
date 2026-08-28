import { opisOdmowy } from '../../komponenty/odmowa';
import { pobierzPlik, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { manifestRepozytorium, metadaneCsv } from './eksporty-biblioteki';
import type { StanBiblioteki } from './stan-biblioteki';

/**
 * Zakładka Archiwum panelu metadanych utrwala i przenosi zasoby: wywozi opis
 * zasobów jako manifest i metadane CSV, a samych bajtów zasobów bez komendy
 * rdzenia wywieźć nie może.
 */
export interface ArchiwumRepozytorium {
  element: HTMLElement;
  /** Przerysowuje wskaźniki liczbowe ze stanu modułu. */
  odswiez(): void;
}

export function utworzArchiwumRepozytorium(stan: StanBiblioteki): ArchiwumRepozytorium {
  const odpowiedz = utworzWierszOdpowiedzi();

  const stanZbioru = document.createElement('p');
  stanZbioru.className = 'dn-pole-opis ml-archiwum__stan';

  const retencja = document.createElement('p');
  retencja.className = 'dn-pole-opis ml-archiwum__retencja';
  retencja.textContent =
    'Retencja: okresu przechowywania okno nie zna. Ustawienie `retencja_historii` należy ' +
    'do okna konfiguracji, a kontrakt nie ma komendy, którą moduł odczytałby je w swoim ' +
    'zasięgu. Wartość domyślna opisana w dokumentacji to przechowywanie bez limitu.';

  const manifest = przycisk('Manifest repozytorium (JSON)', 'dn-btn dn-btn--sm dn-btn--zarys');
  manifest.dataset['czynnosc'] = 'manifest';
  manifest.addEventListener('click', () => wywiez('manifest'));

  const metadane = przycisk('Metadane repozytorium (CSV)', 'dn-btn dn-btn--sm dn-btn--zarys');
  metadane.dataset['czynnosc'] = 'metadane-csv';
  metadane.addEventListener('click', () => wywiez('csv'));

  function wywiez(postac: 'manifest' | 'csv'): void {
    const pliki = stan.pliki();
    if (pliki.length === 0) {
      odpowiedz.pokaz(
        'Wykaz plików jest pusty — wywóz nie miałby ani jednej pozycji. Odśwież wykaz ' +
          'w Library Explorer.',
        false,
      );
      return;
    }
    if (postac === 'manifest') {
      pobierzPlik('manifest-repozytorium.json', manifestRepozytorium(pliki), 'application/json');
    } else {
      pobierzPlik('metadane-repozytorium.csv', metadaneCsv(pliki), 'text/csv');
    }
    odpowiedz.pokaz(
      `Wywieziono opis zasobów z wykazu: ${pliki.length}. Wytwór zawiera metrykę, etykiety ` +
        'i sumy kontrolne — treści plików w nim nie ma, bo kontrakt nie niesie komendy ' +
        'pobierającej bajty zasobu biblioteki.',
      true,
    );
  }

  /** Zasoby objęte czynnością: zaznaczenie, a przy jego braku plik czynny. */
  function objete(): string[] {
    const zaznaczone = stan.zaznaczone();
    if (zaznaczone.length > 0) return [...zaznaczone];
    const czynny = stan.czynny();
    return czynny === null ? [] : [czynny.id];
  }

  /** Utrwalenie archiwalne w jednej z trzech postaci. */
  async function utrwal(postac: 'pdfa' | 'bagit' | 'premis'): Promise<void> {
    const zasoby = objete();
    if (zasoby.length === 0) {
      odpowiedz.pokaz(
        'Utrwalenie nie ma czego objąć — zaznacz zasoby w Library Explorer albo wskaż plik.',
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Rdzeń utrwala zasoby w postaci ${postac}…`, true);
    const wynik = await stan.zrodlo.utrwal(zasoby, postac);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Utrwalenie', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const raporty = wynik.wynik.results.map((pozycja) => pozycja.report).join(' | ');
    odpowiedz.pokaz(
      `Utrwalenie ${postac}: poprawnych ${wynik.wynik.validCount}, niepoprawnych ` +
        `${wynik.wynik.invalidCount}. Zapis walidacji: ${raporty}`,
      wynik.wynik.invalidCount === 0,
    );
    await stan.odczytaj('', '');
  }

  /** Paczka migracyjna albo migawka całego repozytorium. */
  async function spakuj(postac: 'migration' | 'snapshot'): Promise<void> {
    const zasoby = postac === 'snapshot' ? [] : objete();
    if (postac === 'migration' && zasoby.length === 0) {
      odpowiedz.pokaz(
        'Paczka migracyjna obejmuje wskazane zasoby — zaznacz je w Library Explorer. ' +
          'Migawka obejmuje całość i wskazania nie potrzebuje.',
        false,
      );
      return;
    }
    odpowiedz.pokaz('Rdzeń składa archiwum wraz z indeksem i manifestem…', true);
    const wynik = await stan.zrodlo.wywiezPaczke(zasoby, postac, 'zip');
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Paczka', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      `Paczka odłożona w repozytorium jako zasób ${wynik.wynik.fileId}: pozycji ` +
        `${wynik.wynik.entries}, rozmiar ${wynik.wynik.sizeBytes} B, suma manifestu ` +
        `${wynik.wynik.manifestChecksum.slice(0, 16)}…`,
      true,
    );
    await stan.odczytaj('', '');
  }

  /** Raport retencji z rdzenia — polityki i zasoby zbliżające się do terminu. */
  async function odczytajRetencje(): Promise<void> {
    const wynik = await stan.zrodlo.wykazRetencji(30);
    if (!wynik.udany || wynik.wynik === undefined) {
      retencja.textContent = opisOdmowy(
        'Raport retencji',
        wynik.blad?.code,
        wynik.blad?.message,
      );
      return;
    }
    const polityki = wynik.wynik.policies;
    const terminy = wynik.wynik.due ?? [];
    retencja.textContent =
      polityki.length === 0
        ? 'Retencja: ani jednej polityki przechowywania. Wartość domyślna to przechowywanie ' +
          'trwałe — nic nie znika samo.'
        : `Retencja: polityk ${polityki.length} (${polityki
            .map((polityka) => `${polityka.scope}: ${polityka.keepDays} dni → ${polityka.action}`)
            .join(', ')}). Zasobów z terminem w 30 dniach: ${terminy.length}.`;
  }

  const pdfa = przycisk('Utrwal PDF/A', 'dn-btn dn-btn--sm dn-btn--zarys');
  pdfa.dataset['czynnosc'] = 'utrwalenie-pdfa';
  pdfa.addEventListener('click', () => void utrwal('pdfa'));

  const bagit = przycisk('Pakiet BagIt', 'dn-btn dn-btn--sm dn-btn--zarys');
  bagit.dataset['czynnosc'] = 'utrwalenie-bagit';
  bagit.addEventListener('click', () => void utrwal('bagit'));

  const premis = przycisk('Metadane PREMIS/METS', 'dn-btn dn-btn--sm dn-btn--zarys');
  premis.dataset['czynnosc'] = 'utrwalenie-premis';
  premis.addEventListener('click', () => void utrwal('premis'));

  const paczka = przycisk('Paczka migracyjna (ZIP)', 'dn-btn dn-btn--sm dn-btn--zarys');
  paczka.dataset['czynnosc'] = 'paczka-migracyjna';
  paczka.addEventListener('click', () => void spakuj('migration'));

  const migawka = przycisk('Migawka repozytorium', 'dn-btn dn-btn--sm dn-btn--zarys');
  migawka.dataset['czynnosc'] = 'migawka';
  migawka.addEventListener('click', () => void spakuj('snapshot'));

  const pasek = document.createElement('div');
  pasek.className = 'ml-archiwum__pasek';
  pasek.append(pdfa, bagit, premis, paczka, migawka, manifest, metadane);

  const naglowek = document.createElement('h4');
  naglowek.className = 'ml-metadane__naglowek';
  naglowek.textContent = 'Utrwalenie i przeniesienie';

  const element = document.createElement('div');
  element.className = 'ml-archiwum';
  element.append(naglowek, stanZbioru, pasek, retencja, odpowiedz.element);

  return {
    element,

    odswiez() {
      void odczytajRetencje();
      const pliki = stan.pliki();
      const bezSumy = pliki.filter((plik) => (plik.checksum ?? '') === '').length;
      const bajty = pliki.reduce((suma, plik) => suma + (plik.sizeBytes ?? 0), 0);
      stanZbioru.textContent =
        `Wykaz: plików ${pliki.length}, łącznie ${bajty} B wg metryki rdzenia. ` +
        `Bez sumy kontrolnej: ${bezSumy} — pozycja bez sumy wejdzie do manifestu z pustym ` +
        'polem, bo manifest powtarza odpowiedź rdzenia, a nie ją uzupełnia.';
    },
  };
}
