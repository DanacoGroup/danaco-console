import {
  StudioOperationScope,
  type StudioContextualOpRequest,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';

/** Uruchomienie operacji kontekstowej, czynność Tools Panelu wyjęta z wytwórni: zakres bierze się ze stanu modułu, nie z kontrolki wyboru. */
export interface ZapleczeNarzedzi {
  stan: StanStudio;
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
}

/** Składa treść żądania operacji kontekstowej na podstawie stanu modułu; zwraca null, gdy brakuje dokumentu albo wyboru akcji. */
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

/** Wysyła operację kontekstową do rdzenia i odkłada jej wynik jako propozycję zmiany dokumentu w stanie modułu. */
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
    // Odmowa idzie też do stanu modułu, inaczej kanwa tekstowa pokazałaby brak propozycji jako pustkę.
    zaplecze.stan.ustawOdmoweOperacji(`Operacja ${idAkcji}: ${powod}`);
    return;
  }
  const tresc = wynik.wynik.resultText ?? '';
  const idPropozycji = wynik.wynik.proposalId ?? '';
  zaplecze.stan.ustawPropozycje({ idPropozycji, tresc, idAkcji });
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz(opiszWynikOperacji(tresc, idPropozycji), true);
}

/** Zdanie opisujące treść, która przyszła z operacji kontekstowej, zamiast ogólnikowego stwierdzenia, że wynik czeka na przyjęcie. */
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
