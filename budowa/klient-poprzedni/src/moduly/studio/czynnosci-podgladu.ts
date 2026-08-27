import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { wydajDokument } from './konwersja-dokumentu';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';
import type { ZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import type { ZrodloPrzekazania } from './zrodlo-przekazania';

/** Eksport i przekazanie do Library, dwie czynności Preview Window wyjęte z wytwórni okna: obie wychodzą do rdzenia i obie wykonują się w całości. */
export interface ZapleczePodgladu {
  stan: StanStudio;
  /** Źródło komend obszaru dokumentów — nośnik zamiany formatu. */
  dokumenty: ZrodloDokumentuStudio;
  /** Pas stanu okna — nośnik ładowania i odmowy. */
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
  /** Przywraca pasowi stan wynikający z danych, gdy komunikat ustępuje. */
  odswiez(): void;
}

export async function eksportujDokument(
  zaplecze: ZapleczePodgladu,
  format: string,
): Promise<void> {
  const { stan, dokumenty, pas, odpowiedz } = zaplecze;
  pas.ladowanie(`Wydanie dokumentu w formacie ${format} w toku…`);
  const udane = await wydajDokument({ stan, dokumenty, odpowiedz }, format);
  if (!udane) {
    // Powód stoi już w wierszu odpowiedzi, więc pas wraca do stanu wynikającego z danych bez powtórzenia.
    pas.gotowe();
    zaplecze.odswiez();
    return;
  }
  pas.gotowe();
  zaplecze.odswiez();
}

/**
 * Przekazanie dokumentu do Library kompletem kontekstu.
 *
 * Komplet niesie identyfikator dokumentu, bo to on jest przedmiotem
 * przekazania; polecenie wyjściowe zostaje puste, bo podgląd nie zleca pracy,
 * tylko przenosi wynik.
 */
export async function przekazDoLibrary(
  zaplecze: ZapleczePodgladu,
  przekazanie: ZrodloPrzekazania,
): Promise<void> {
  const { stan, pas, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  if (dokument === null) {
    odpowiedz.pokaz('Przekazanie dotyczy dokumentu — wczytaj go najpierw w Studio Editorze.', false);
    return;
  }
  pas.ladowanie('Przekazywanie kontekstu do Library…');
  odpowiedz.pokaz('Przekazywanie kontekstu do Library…', true);
  const wynik = await przekazanie.doLibrary(stan.idOkna(), { documentIds: [dokument.id] });
  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowy('Przekazanie do Library', wynik.blad?.code, wynik.blad?.message);
    pas.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  pas.gotowe();
  zaplecze.odswiez();
  odpowiedz.pokaz(
    `Kontekst przeniesiony do okna ${wynik.wynik.window.id} modułu Library ` +
      `(przeniesiono: ${String(wynik.wynik.transferred)}).`,
    true,
  );
}
