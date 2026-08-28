import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { StanBiblioteki } from './stan-biblioteki';
import { utworzSterModulu } from './ster-modulu';
import {
  nadajEtykieteZbiorczo,
  przeniesZasoby,
  przypiszZaznaczenie,
  type OdpowiedzZapisu,
} from './zapisy-zbiorcze';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Czynności zbiorcze Library Explorera obejmują osiem operacji z drogą do
 * rdzenia: etykietę, przypisanie do kolekcji, przeniesienie, archiwizację,
 * usunięcie trwałe i udostępnienie odnośnikiem; usunięcie pyta osobno, bo jest
 * jedyną nieodwracalną.
 */
export interface CzynnosciZbiorcze {
  element: HTMLElement;
  /** Przerysowuje podpowiedzi kolekcji i licznik zaznaczenia. */
  odswiez(): void;
}

export function utworzCzynnosciZbiorcze(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
): CzynnosciZbiorcze {
  const odpowiedz = utworzWierszOdpowiedzi();

  const etykieta = poleTekstowe({
    etykieta: 'Etykieta zbiorcza',
    podpowiedz: 'np. pisma-procesowe',
    opis: 'Etykieta dochodzi do etykiet już nadanych; nie zastępuje ich.',
  });
  const kolekcja = poleWyboru({ etykieta: 'Kolekcja docelowa' }, []);
  const modul = utworzSterModulu(
    stan,
    'Katalog modułów pochodzi z komendy module.list rdzenia. Ta sama nastawa stoi w oknie ' +
      'File Preview — zmiana tutaj jest widoczna i tam.',
  );

  const sciezka = poleTekstowe({
    etykieta: 'Ścieżka docelowa w repozytorium',
    podpowiedz: 'np. Klienci/Klient X',
    opis: 'Przeniesienie zmienia miejsce w porządku repozytorium; treści ani wersji nie rusza.',
  });

  const licznik = document.createElement('p');
  licznik.className = 'dn-pole-opis ml-czynnosci__licznik';

  /** Zasoby wskazane do czynności nieodwracalnej — zapamiętane między kliknięciami. */
  let potwierdzaneUsuniecie: string[] = [];

  /** Wykonuje czynność rdzenia na zaznaczeniu i odświeża wykaz. */
  async function naZaznaczeniu(
    zapowiedz: string,
    praca: () => Promise<{ udany: boolean; blad?: { code?: string; message?: string } | undefined; zdanie: string }>,
  ): Promise<void> {
    if (stan.zaznaczone().length === 0) {
      odpowiedz.pokaz('Nie zaznaczono ani jednego pliku — czynność nie ma czego objąć.', false);
      return;
    }
    odpowiedz.pokaz(zapowiedz, true);
    const wynik = await praca();
    odpowiedz.pokaz(
      wynik.udany
        ? wynik.zdanie
        : opisOdmowy('Czynność zbiorcza', wynik.blad?.code, wynik.blad?.message),
      wynik.udany,
    );
    if (wynik.udany) await stan.odczytaj('', '');
  }

  /** Pliki zaznaczone w wykazie — wspólne wejście trzech czynności. */
  function wybrane() {
    const zaznaczone = stan.zaznaczone();
    return stan.pliki().filter((plik) => zaznaczone.includes(plik.id));
  }

  /** Wykonuje czynność i pokazuje jej jedno zdanie; zapowiedź trwania włącznie. */
  async function wykonaj(zapowiedz: string, praca: Promise<OdpowiedzZapisu>): Promise<void> {
    odpowiedz.pokaz(zapowiedz, true);
    const wynik = await praca;
    odpowiedz.pokaz(wynik.tresc, wynik.powodzenie);
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-czynnosci__pasek';
  pasek.append(
    przyciskCzynnosci('Tag', 'tag', () =>
      wykonaj(
        'Nadawanie etykiety zaznaczonym plikom…',
        nadajEtykieteZbiorczo(stan, etykieta.kontrolka.value, wybrane()),
      ),
    ),
    przyciskCzynnosci('Do kolekcji', 'kolekcja', () =>
      wykonaj(
        'Przypisywanie zaznaczenia do kolekcji…',
        przypiszZaznaczenie(stan, kolekcja.kontrolka.value, wybrane()),
      ),
    ),
    przyciskCzynnosci('Otwórz w module źródłowym', 'otworz', () =>
      wykonaj(
        'Przenoszenie zaznaczenia do modułu docelowego…',
        przeniesZasoby(stan, otoczenie, stan.modulDocelowy(), wybrane()),
      ),
    ),
    przyciskCzynnosci('Przenieś w repozytorium', 'przeniesienie', () =>
      naZaznaczeniu('Przenoszenie zaznaczenia pod wskazaną ścieżkę…', async () => {
        const wynik = await stan.zrodlo.przenies(stan.zaznaczone(), sciezka.kontrolka.value);
        return {
          udany: wynik.udany && wynik.wynik !== undefined,
          blad: wynik.blad,
          zdanie: `Przeniesiono zasobów: ${wynik.wynik?.movedCount ?? 0} pod ścieżkę ` +
            `${sciezka.kontrolka.value}.`,
        };
      }),
    ),
    przyciskCzynnosci('Archiwizuj', 'archiwizacja', () =>
      naZaznaczeniu('Przenoszenie zaznaczenia do archiwum…', async () => {
        const wynik = await stan.zrodlo.zarchiwizuj(stan.zaznaczone(), '');
        return {
          udany: wynik.udany && wynik.wynik !== undefined,
          blad: wynik.blad,
          zdanie: `Zarchiwizowano zasobów: ${wynik.wynik?.archivedCount ?? 0}. Czynność jest ` +
            'odwracalna przyciskiem „Przywróć z archiwum".',
        };
      }),
    ),
    przyciskCzynnosci('Przywróć z archiwum', 'przywrocenie', () =>
      naZaznaczeniu('Przywracanie zaznaczenia z archiwum…', async () => {
        const wynik = await stan.zrodlo.przywrocZArchiwum(stan.zaznaczone());
        return {
          udany: wynik.udany && wynik.wynik !== undefined,
          blad: wynik.blad,
          zdanie: `Przywrócono zasobów: ${wynik.wynik?.restoredCount ?? 0}.`,
        };
      }),
    ),
    przyciskCzynnosci('Udostępnij odnośnikiem', 'udostepnienie', async () => {
      const zaznaczone = stan.zaznaczone();
      if (zaznaczone.length !== 1) {
        odpowiedz.pokaz(
          'Odnośnik wystawia się do jednego zasobu — zaznacz dokładnie jeden plik.',
          false,
        );
        return;
      }
      odpowiedz.pokaz('Rdzeń wystawia odnośnik wraz z tokenem dostępu…', true);
      const wynik = await stan.zrodlo.wystawUdostepnienie('file', zaznaczone[0] ?? '');
      if (!wynik.udany || wynik.wynik === undefined) {
        odpowiedz.pokaz(
          opisOdmowy('Udostępnienie', wynik.blad?.code, wynik.blad?.message),
          false,
        );
        return;
      }
      odpowiedz.pokaz(
        `Odnośnik wystawiony: ${wynik.wynik.url}. Token jest jawny — przekazanie go komuś ` +
          'przekazuje dostęp do zasobu; odwołać go można komendą library.share.revoke.',
        true,
      );
    }),
    przyciskCzynnosci('Usuń trwale', 'usuniecie', async () => {
      const zaznaczone = [...stan.zaznaczone()];
      if (zaznaczone.length === 0) {
        odpowiedz.pokaz('Nie zaznaczono ani jednego pliku.', false);
        return;
      }
      const takieSame =
        potwierdzaneUsuniecie.length === zaznaczone.length &&
        potwierdzaneUsuniecie.every((kod) => zaznaczone.includes(kod));
      if (!takieSame) {
        potwierdzaneUsuniecie = zaznaczone;
        odpowiedz.pokaz(
          `Usunięcie trwałe zdejmie ${zaznaczone.length} zasobów wraz z ich wersjami ` +
            'i treścią. Czynność jest nieodwracalna — naciśnij ponownie, żeby ją wykonać.',
          false,
        );
        return;
      }
      potwierdzaneUsuniecie = [];
      odpowiedz.pokaz('Usuwanie trwałe zaznaczonych zasobów…', true);
      const wynik = await stan.zrodlo.usunTrwale(zaznaczone);
      if (!wynik.udany || wynik.wynik === undefined) {
        odpowiedz.pokaz(opisOdmowy('Usunięcie', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      odpowiedz.pokaz(
        `Usunięto trwale zasobów: ${wynik.wynik.deletedCount}, wersji: ` +
          `${wynik.wynik.deletedVersions}.`,
        true,
      );
      await stan.odczytaj('', '');
    }),
  );

  const element = document.createElement('div');
  element.className = 'ml-czynnosci';
  element.append(
    licznik,
    etykieta.element,
    kolekcja.element,
    sciezka.element,
    modul.element,
    pasek,
    odpowiedz.element,
  );

  return {
    element,
    odswiez() {
      licznik.textContent = `Zaznaczono plików: ${stan.zaznaczone().length}.`;
      ustawPozycje(
        kolekcja.kontrolka,
        stan.kolekcje().map((kod) => ({ wartosc: kod, etykieta: kod })),
      );
      modul.odswiez();
    },
  };
}

/** Przycisk czynności zbiorczej wraz ze znacznikiem do sprawdzianu jednoznacznie wskazującym rodzaj czynności. */
function przyciskCzynnosci(
  etykieta: string,
  kod: string,
  czynnosc: () => Promise<void>,
): HTMLElement {
  const element = przycisk(etykieta, 'dn-btn dn-btn--sm dn-btn--zarys');
  element.dataset['czynnosc'] = kod;
  element.addEventListener('click', () => void czynnosc());
  return element;
}
