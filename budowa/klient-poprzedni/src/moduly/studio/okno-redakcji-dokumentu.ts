import type { DesignAsset } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { CZYNNOSCI_REDAKCJI, opiszSkutekRedakcji } from './czynnosci-redakcji';
import type { CzynnoscWarsztatu } from './czynnosci-warsztatu';
import { utworzKontrolke, type KontrolkaPola } from './okno-warsztatu-dokumentu';
import { utworzOknoStudio } from './okno-studio';
import type { StanStudio } from './stan-studio';
import type { ZrodloWarsztatuDokumentu } from './zrodlo-warsztatu-dokumentu';

/** Redakcja dokumentu, okno piętnastu czynności prowadzących dokument Studia poza edytor, na osobnym materiale niż warsztat. */
export interface OknoRedakcjiDokumentu {
  element: HTMLElement;
  /** Wczytuje wykaz zasobów okna — potrzebny polom osadzenia i porównania. */
  wczytaj(): Promise<void>;
  odswiez(): void;
}

export function utworzOknoRedakcjiDokumentu(
  stan: StanStudio,
  zrodlo: ZrodloWarsztatuDokumentu,
): OknoRedakcjiDokumentu {
  const rama = utworzOknoStudio({
    kod: 'studio.document-redaction',
    tytul: 'Redakcja dokumentu',
    rola: 'pomocnicze',
    objasnienie:
      'Czynności prowadzące dokument poza edytor: gałęzie i ich scalanie, wydanie ' +
      'historii i paczki redakcyjnej, raport zmian, podgląd układu, różnica wizualna, ' +
      'porównanie ze źródłem, wyszukiwanie znaczeniowe, osadzenie zasobu, wsad, ' +
      'wczytanie strony sieciowej. Archiwa składa Go, wyrys idzie biblioteką ' +
      'wkompilowaną w rdzeń — bez ani jednego programu zewnętrznego.',
  });

  const wybor = poleWyboru(
    {
      etykieta: 'Czynność',
      opis: 'Wybór przestawia pola poniżej — każda czynność bierze inne wskazania.',
    },
    CZYNNOSCI_REDAKCJI.map((czynnosc) => ({
      wartosc: czynnosc.komenda,
      etykieta: czynnosc.nazwa,
    })),
  );

  const objasnienieCzynnosci = document.createElement('p');
  objasnienieCzynnosci.className = 'dn-pole-opis ms-wskaznik';

  const formularz = document.createElement('div');
  formularz.className = 'ms-warsztat__pola';

  const wykonanie = przycisk('Wykonaj czynność', 'dn-btn dn-btn--sygnal');
  wykonanie.dataset['czynnosc'] = 'wykonaj-redakcje';
  const odpowiedz = utworzWierszOdpowiedzi();

  rama.pasek.append(wykonanie);
  rama.stan.tresc.append(wybor.element, objasnienieCzynnosci, formularz, odpowiedz.element);

  let zasoby: DesignAsset[] = [];
  let kontrolki = new Map<string, KontrolkaPola>();

  /** Czynność wskazana na liście; pierwsza, gdy lista dopiero się zakłada. */
  function czynnoscBiezaca(): CzynnoscWarsztatu {
    const wybrana = CZYNNOSCI_REDAKCJI.find((czynnosc) => czynnosc.komenda === wybor.kontrolka.value);
    return wybrana ?? (CZYNNOSCI_REDAKCJI[0] as CzynnoscWarsztatu);
  }

  /** Przebudowuje formularz pod czynność wskazaną. */
  function przestawFormularz(): void {
    const czynnosc = czynnoscBiezaca();
    objasnienieCzynnosci.textContent = czynnosc.opis;
    kontrolki = new Map();
    const wiersze: HTMLElement[] = [];
    for (const pole of czynnosc.pola) {
      const kontrolka = utworzKontrolke(pole, zasoby);
      kontrolki.set(pole.kod, kontrolka);
      wiersze.push(kontrolka.element);
    }
    formularz.replaceChildren(...wiersze);
  }

  /** Zbiera wartości formularza w postaci, którą czyta katalog czynności. */
  function zbierzWartosci(): Record<string, string> {
    const wartosci: Record<string, string> = {};
    for (const [kod, kontrolka] of kontrolki) wartosci[kod] = kontrolka.odczytaj();
    return wartosci;
  }

  async function wykonaj(): Promise<void> {
    const czynnosc = czynnoscBiezaca();
    const zlozenie = czynnosc.zloz(zbierzWartosci(), {
      idOkna: stan.idOkna(),
      idDokumentu: stan.dokument()?.id ?? null,
    });
    if ('odmowa' in zlozenie) {
      odpowiedz.pokaz(zlozenie.odmowa, false);
      return;
    }

    rama.stan.ladowanie(`${czynnosc.nazwa} — czynność w toku…`);
    const wynik = await zrodlo.wykonaj(czynnosc.komenda, zlozenie.zadanie);
    if (!wynik.udany) {
      // Odmowa rdzenia jest tu wynikiem pełnoprawnym: skanowanie odmawia z zasady, ma to być przeczytane.
      odpowiedz.pokaz(opisOdmowy(czynnosc.nazwa, wynik.blad?.code, wynik.blad?.message), false);
      rama.stan.gotowe();
      return;
    }
    odpowiedz.pokaz(`${czynnosc.nazwa}: ${opiszSkutekRedakcji(wynik.wynik)}`, true);
    rama.stan.gotowe();
    // Wydanie, wyrys i wczytanie strony zakładają nowe zasoby, więc wykaz zestarzał się w tej chwili.
    await wczytaj();
  }

  async function wczytaj(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') return;
    const wynik = await zrodlo.zasoby(idOkna);
    // Odmowa odczytu zasobów nie gasi okna: zasoby wypełniają dwa pola, reszta czynności działa bez nich.
    zasoby = wynik.udany && wynik.wynik !== undefined ? wynik.wynik : [];
    przestawFormularz();
    odswiez();
  }

  function odswiez(): void {
    if (rama.stan.faza() === 'ladowanie') return;
    rama.stan.gotowe();
  }

  wybor.kontrolka.addEventListener('change', () => {
    odpowiedz.wyczysc();
    przestawFormularz();
  });
  wykonanie.addEventListener('click', () => void wykonaj());

  przestawFormularz();
  rama.stan.gotowe();

  return { element: rama.element, wczytaj, odswiez };
}
