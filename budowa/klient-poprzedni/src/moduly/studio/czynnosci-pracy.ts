import { StudioOperationScope, type StudioTrackedChange } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';
import { ladunekOperacji, zdanieDlaModelu, type NastawySuwakow } from './suwaki-koncepcyjne';
import type { ZrodloPracyStudio } from './zrodlo-pracy-studio';

/** Czynności okna pracy z dokumentem, które rozmawiają z rdzeniem: operacja z poleceniem własnym, decyzja o zmianach, śledzenie i komentarze. */
export interface ZapleczePracy {
  stan: StanStudio;
  praca: ZrodloPracyStudio;
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
  /** Wywoływane po każdej zmianie, którą rdzeń potwierdził. */
  poZmianie(): void;
}

/** Zleca operację kontekstową wraz z poleceniem operatora i nastawami suwaków, biorąc zakres skuteczny ze stanu modułu. */
export async function zlecOperacje(
  zaplecze: ZapleczePracy,
  idAkcji: string,
  polecenie: string,
  nastawy: NastawySuwakow,
): Promise<void> {
  const { stan, pas, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  if (dokument === null) {
    odpowiedz.pokaz(
      'Operacja wymaga dokumentu wczytanego z rdzenia — komenda studio.contextual.op przyjmuje ' +
        'jego identyfikator jako pole obowiązkowe.',
      false,
    );
    return;
  }
  if (idAkcji === '') {
    odpowiedz.pokaz('Wskaż operację — żądanie bez pozycji rejestru akcji nie ma czego zlecić.', false);
    return;
  }

  const wybrany = stan.zaznaczenie();
  const naZaznaczeniu = stan.zakresSkuteczny() === StudioOperationScope.Selection;
  const zdanie = zdanieDlaModelu(nastawy, polecenie);
  pas.ladowanie(`Operacja ${idAkcji} w toku…`);

  const wynik = await stan.zrodlo.operacja({
    windowId: stan.idOkna(),
    documentId: dokument.id,
    actionId: idAkcji,
    scope: naZaznaczeniu ? StudioOperationScope.Selection : StudioOperationScope.Document,
    ...(naZaznaczeniu && wybrany !== null
      ? { selectionStart: wybrany.poczatek, selectionEnd: wybrany.koniec }
      : {}),
    params: ladunekOperacji(nastawy, polecenie) as never,
  });

  if (!wynik.udany || wynik.wynik === undefined) {
    const powod = opisOdmowy('Operacja kontekstowa', wynik.blad?.code, wynik.blad?.message);
    pas.blad(powod);
    stan.ustawOdmoweOperacji(`Operacja ${idAkcji}: ${powod}`);
    return;
  }
  const tresc = wynik.wynik.resultText ?? '';
  const idPropozycji = wynik.wynik.proposalId ?? '';
  stan.ustawPropozycje({ idPropozycji, tresc, idAkcji });
  pas.gotowe();
  odpowiedz.pokaz(
    `Rdzeń wykonał operację ${idAkcji}${zdanie === '' ? '' : ` (${zdanie})`} i wpisał wynik ` +
      `(${tresc.length} znaków) do treści dokumentu jako zmianę autora „model". Zmiana czeka ` +
      'na decyzję w miejscu, w którym stoi; propozycja ' +
      (idPropozycji === '' ? 'nie została odłożona' : idPropozycji) +
      ' zostaje do porównania stron.',
    true,
  );
  zaplecze.poZmianie();
}

/** Odczytuje zmiany śledzone dokumentu czynnego z rdzenia; brak wczytanego dokumentu oddaje wykaz pusty zamiast odmowy. */
export async function odczytajZmiany(
  zaplecze: ZapleczePracy,
): Promise<readonly StudioTrackedChange[]> {
  const dokument = zaplecze.stan.dokument();
  if (dokument === null) return [];
  const wynik = await zaplecze.praca.zmiany(dokument.id, true);
  if (!wynik.udany || wynik.wynik === undefined) return [];
  return wynik.wynik.changes;
}

/** Rozstrzyga wskazane zmiany śledzone jako przyjęte albo odrzucone i wciąga do stanu dokument oddany przez rdzeń. */
export async function rozstrzygnijZmiany(
  zaplecze: ZapleczePracy,
  kody: readonly string[],
  przyjmij: boolean,
): Promise<void> {
  const { stan, praca, pas, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  if (dokument === null || kody.length === 0) {
    odpowiedz.pokaz(
      'Decyzja dotyczy zmian wskazanych — komenda studio.tracking.decide przyjmuje ich ' +
        'identyfikatory i wykazu pustego nie przyjmie.',
      false,
    );
    return;
  }
  pas.ladowanie(przyjmij ? 'Przyjmowanie zmian…' : 'Odrzucanie zmian…');
  const wynik = await praca.rozstrzygnijZmiany(dokument.id, kody, przyjmij);
  if (!wynik.udany || wynik.wynik === undefined) {
    pas.blad(opisOdmowy('Decyzja o zmianach', wynik.blad?.code, wynik.blad?.message));
    return;
  }
  stan.wchlon(wynik.wynik.document);
  pas.gotowe();
  odpowiedz.pokaz(
    `${przyjmij ? 'Przyjęto' : 'Odrzucono'} zmian: ${wynik.wynik.decided}. Rdzeń zapisał treść ` +
      'i założył wersję w repozytorium sesji.',
    true,
  );
  zaplecze.poZmianie();
}

/** Przestawia śledzenie zmian dokumentu w rdzeniu i oddaje wywołującemu stan śledzenia odczytany z odpowiedzi rdzenia. */
export async function przestawSledzenie(
  zaplecze: ZapleczePracy,
  czynne: boolean,
): Promise<boolean> {
  const dokument = zaplecze.stan.dokument();
  if (dokument === null) {
    zaplecze.odpowiedz.pokaz(
      'Śledzenie przestawia się na dokumencie — wczytaj go najpierw.',
      false,
    );
    return false;
  }
  const wynik = await zaplecze.praca.ustawSledzenie(dokument.id, czynne);
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.odpowiedz.pokaz(
      opisOdmowy('Przestawienie śledzenia', wynik.blad?.code, wynik.blad?.message),
      false,
    );
    return false;
  }
  zaplecze.odpowiedz.pokaz(
    wynik.wynik.enabled
      ? 'Śledzenie zmian włączone — zapisy odkładają wstawienia i usunięcia do decyzji.'
      : 'Śledzenie zmian wyłączone. Wynik operacji modelu wchodzi do treści jako zmiana ' +
          'oznaczona niezależnie od tej nastawy — tak stanowi droga operacji kontekstowej.',
    true,
  );
  return wynik.wynik.enabled;
}

/** Rozstrzyga propozycję zmiany w rdzeniu, przyjmując ją w całości albo tylko wskazanymi fragmentami treści. */
export async function rozstrzygnijPropozycje(
  zaplecze: ZapleczePracy,
  przyjmij: boolean,
  fragmenty: readonly number[],
): Promise<void> {
  const { stan, praca, pas, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  const propozycja = stan.propozycja();
  if (propozycja === null) {
    odpowiedz.pokaz(
      'Nie ma wyniku operacji do decyzji — zleć operację wierszem polecenia przy kursorze, ' +
        'suwakiem na wstążce albo z Tools Panelu.',
      false,
    );
    return;
  }
  if (dokument === null || propozycja.idPropozycji === '') {
    // Propozycja bez odwołania w rdzeniu istnieje tylko w tym oknie, więc decyzja zapada tutaj.
    if (przyjmij) stan.przyjmijPropozycje();
    else stan.ustawPropozycje(null);
    odpowiedz.pokaz(
      przyjmij
        ? 'Wynik przyjęty po stronie okna — rdzeń nie odłożył propozycji, więc nie ma czego ' +
            'rozstrzygać u niego. Zapis utrwali treść i założy wersję.'
        : 'Wynik odrzucony; treść dokumentu została bez zmiany.',
      true,
    );
    zaplecze.poZmianie();
    return;
  }
  pas.ladowanie('Decyzja o propozycji w toku…');
  const wynik = await praca.rozstrzygnijPropozycje(
    dokument.id,
    propozycja.idPropozycji,
    przyjmij,
    fragmenty,
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    pas.blad(opisOdmowy('Decyzja o propozycji', wynik.blad?.code, wynik.blad?.message));
    return;
  }
  stan.wchlon(wynik.wynik.document);
  stan.ustawPropozycje(null);
  pas.gotowe();
  odpowiedz.pokaz(
    przyjmij
      ? `Rdzeń przyjął ${fragmenty.length === 0 ? 'całą propozycję' : `fragmenty ${fragmenty.join(', ')}`}` +
          ' i złożył treść dokumentu z wybranych stron różnicy.'
      : 'Rdzeń odrzucił propozycję; treść dokumentu została bez zmiany.',
    true,
  );
  zaplecze.poZmianie();
}

/** Zakłada nowy komentarz przypięty do zaznaczenia w dokumencie albo odpowiedź w istniejącym wątku komentarzy. */
export async function dodajKomentarz(
  zaplecze: ZapleczePracy,
  tresc: string,
  zakres: { poczatek: number; koniec: number } | null,
  idWatku: string,
): Promise<void> {
  const { stan, praca, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  if (dokument === null) {
    odpowiedz.pokaz('Komentarz dotyczy dokumentu — wczytaj go najpierw.', false);
    return;
  }
  const wynik = await praca.dodajKomentarz(dokument.id, tresc, zakres, idWatku);
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(opisOdmowy('Zapis komentarza', wynik.blad?.code, wynik.blad?.message), false);
    return;
  }
  odpowiedz.pokaz(
    idWatku === ''
      ? `Komentarz ${wynik.wynik.comment.id} przypięty${zakres === null ? ' do dokumentu' : ' do zaznaczenia'}.`
      : `Odpowiedź dopisana do wątku ${idWatku}.`,
    true,
  );
  zaplecze.poZmianie();
}
