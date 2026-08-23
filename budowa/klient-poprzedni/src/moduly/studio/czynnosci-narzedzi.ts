import {
  StudioOperationScope,
  type StudioContextualOpRequest,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';

/**
 * Uruchomienie operacji kontekstowej — czynność Tools Panelu wyjęta z wytwórni.
 *
 * Zakres bierze się ze stanu modułu, a nie z kontrolki: zakres skuteczny liczy
 * `stan-studio` i biorą go stamtąd wszyscy trzej odbiorcy — wskaźnik panelu,
 * pasek zaznaczenia edytora i to żądanie. Druga kopia nastawy, trzymana w `value`
 * listy wyboru panelu, rozjeżdżała się ze stanem.
 *
 * Wybór „Zaznaczenie" bez zaznaczenia nie jest blokowany — schodzi na cały
 * dokument, a panel pisze o tym we wskaźniku zakresu. Żądanie z zakresem
 * `selection` i bez granic zaznaczenia dostałoby odmowę walidacji.
 */
export interface ZapleczeNarzedzi {
  stan: StanStudio;
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
}

/** Składa treść żądania operacji; `null`, gdy brakuje dokumentu albo wyboru. */
export function zlozZadanieOperacji(
  stan: StanStudio,
  idAkcji: string,
): StudioContextualOpRequest | null {
  const dokument = stan.dokument();
  if (dokument === null || idAkcji === '') return null;
  const wybrany = stan.zaznaczenie();
  const naZaznaczeniu = stan.zakresSkuteczny() === StudioOperationScope.Selection;
  const tresc: StudioContextualOpRequest = {
    windowId: stan.idOkna(),
    documentId: dokument.id,
    actionId: idAkcji,
    scope: naZaznaczeniu ? StudioOperationScope.Selection : StudioOperationScope.Document,
  };
  if (naZaznaczeniu && wybrany !== null) {
    tresc.selectionStart = wybrany.poczatek;
    tresc.selectionEnd = wybrany.koniec;
  }
  return tresc;
}

/** Wysyła operację do rdzenia i odkłada jej wynik jako propozycję zmiany. */
export async function uruchomOperacje(
  zaplecze: ZapleczeNarzedzi,
  zadanie: StudioContextualOpRequest | null,
  idAkcji: string,
): Promise<void> {
  if (zadanie === null) {
    zaplecze.odpowiedz.pokaz(BRAK_WARUNKOW, false);
    return;
  }
  zaplecze.pas.ladowanie(`Operacja ${idAkcji} w toku…`);
  const wynik = await zaplecze.stan.zrodlo.operacja(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    const powod = opisOdmowy('Operacja kontekstowa', wynik.blad?.code, wynik.blad?.message);
    zaplecze.pas.blad(powod);
    // Odmowa idzie także do stanu modułu, nie tylko do pasa tego okna. Kanwa
    // tekstowa stoi obok i bez tego pokazałaby brak propozycji jako „operacji
    // jeszcze nie było" — czyli pustkę w miejscu odmowy rdzenia.
    zaplecze.stan.ustawOdmoweOperacji(`Operacja ${idAkcji}: ${powod}`);
    return;
  }
  const tresc = wynik.wynik.resultText ?? '';
  const idPropozycji = wynik.wynik.proposalId ?? '';
  zaplecze.stan.ustawPropozycje({ idPropozycji, tresc, idAkcji });
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz(opiszWynikOperacji(tresc, idPropozycji), true);
}

/**
 * Zdanie o treści, która przyszła, a nie o tym, że „wynik czeka".
 *
 * `studio.contextual.op` wraca odpowiedzią pomyślną także wtedy, gdy kanał modelu
 * oddał zamiast wyniku swój własny komunikat. Rdzeń tego nie odróżnia, bo dostał
 * fragmenty treści, a klient nie ma po czym, bo to zwykły tekst. Zamiast zgadywać
 * po treści zdanie mówi, ile jej przyszło i że przyjęcie wyniku podmieni nią
 * dokument.
 */
function opiszWynikOperacji(tresc: string, idPropozycji: string): string {
  if (tresc !== '') {
    return (
      `Rdzeń oddał treść wyniku (${tresc.length} znaków) — przeczytaj ją w Diff/Grep Panelu ` +
      'przed przyjęciem, bo przyjęcie podmieni nią treść dokumentu.'
    );
  }
  if (idPropozycji !== '') {
    return (
      `Rdzeń oddał samo odwołanie do propozycji ${idPropozycji}, bez treści wyniku. ` +
      'Porównaj ją w Diff/Grep Panelu — przyjęcie bez treści zostawi dokument bez zmiany.'
    );
  }
  return 'Rdzeń nie oddał ani treści wyniku, ani odwołania do propozycji — nie ma czego przyjąć.';
}

const BRAK_WARUNKOW =
  'Operacja wymaga wczytanego dokumentu i wskazanej pozycji wykazu — komenda studio.contextual.op ' +
  'przyjmuje oba pola jako obowiązkowe.';
