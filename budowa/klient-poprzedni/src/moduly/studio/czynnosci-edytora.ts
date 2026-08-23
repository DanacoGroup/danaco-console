import type { StudioDocumentOpenRequest } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';

/**
 * Cztery czynności Studio Editora wyjęte z jego wytwórni.
 *
 * Widok odpowiada za układ, stany i podpięcie zdarzeń; ten plik za to, co dzieje
 * się po naciśnięciu — rozmowę z rdzeniem i skutek dla stanu modułu. Zmiana
 * kształtu żądania nie zmusza więc do czytania układu okna.
 *
 * Odmowa zostaje w pasie stanu, powodzenie w wierszu odpowiedzi: komunikat błędu
 * ma być trwały w układzie, a potwierdzenie czynności ustępuje następnemu
 * naciśnięciu.
 */
export interface ZapleczeEdytora {
  stan: StanStudio;
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
}

/** Wczytuje dokument; `zadanie` puste znaczy formularz bez wskazania. */
export async function wczytajDokument(
  zaplecze: ZapleczeEdytora,
  zadanie: StudioDocumentOpenRequest | null,
  brak: string,
): Promise<void> {
  if (zadanie === null) {
    zaplecze.odpowiedz.pokaz(brak, false);
    return;
  }
  zaplecze.pas.ladowanie('Wczytywanie dokumentu z rdzenia…');
  const wynik = await zaplecze.stan.zrodlo.otworz(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.pas.blad(opisOdmowy('Wczytanie dokumentu', wynik.blad?.code, wynik.blad?.message));
    return;
  }
  const dokument = wynik.wynik.document;
  zaplecze.stan.wchlon(dokument);
  zaplecze.pas.gotowe();
  const skutek = opiszSkutekWczytania(zadanie, dokument.content ?? '');
  zaplecze.odpowiedz.pokaz(`Dokument ${dokument.id}: ${skutek}`, skutek === TRESC_PRZYSZLA);
}

/**
 * Zdanie o tym, co rdzeń oddał, a nie o tym, o co go poproszono.
 *
 * `studio.document.open` ze wskazaniem `path` oddaje dokument bez pola `content`,
 * bo rdzeń nie ma dostępu do systemu plików Operatora; ze wskazaniem
 * `libraryFileId` oddaje dokument z zapamiętanym odwołaniem i również bez treści,
 * bo doczytanie jej zrobiłoby z adaptera Studio klienta cudzego modułu. Dokument
 * zakłada się więc pusty, a treść dostarcza dopiero pierwszy `studio.document.save`.
 */
const TRESC_PRZYSZLA = 'treść przyszła z rdzenia.';

function opiszSkutekWczytania(zadanie: StudioDocumentOpenRequest, tresc: string): string {
  if (tresc !== '') return TRESC_PRZYSZLA;
  if (zadanie.path !== undefined && zadanie.path !== '') {
    return (
      'rdzeń założył dokument PUSTY i ścieżki na urządzeniu NIE PRZECZYTAŁ — nie ma dostępu ' +
      'do systemu plików Operatora. Treść wklej albo wpisz tutaj; pierwszy zapis ją utrwali.'
    );
  }
  if (zadanie.libraryFileId !== undefined && zadanie.libraryFileId !== '') {
    return (
      'rdzeń zapamiętał odwołanie do pliku Library, ale jego treści NIE DOCZYTAŁ — dokument ' +
      'jest pusty. Treść wklej albo wpisz tutaj; pierwszy zapis ją utrwali.'
    );
  }
  return 'rdzeń oddał dokument bez treści.';
}

/** Zapisuje treść roboczą i zakłada wersję w repozytorium sesji. */
export async function zapiszDokument(zaplecze: ZapleczeEdytora): Promise<void> {
  const dokument = zaplecze.stan.dokument();
  if (dokument === null) {
    zaplecze.odpowiedz.pokaz(
      'Zapis wymaga dokumentu wczytanego z rdzenia — komenda studio.document.save przyjmuje jego identyfikator.',
      false,
    );
    return;
  }
  zaplecze.pas.ladowanie('Zapis dokumentu i zakładanie wersji…');
  const wynik = await zaplecze.stan.zrodlo.zapisz({
    documentId: dokument.id,
    content: zaplecze.stan.trescRobocza(),
    createVersion: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.pas.blad(opisOdmowy('Zapis dokumentu', wynik.blad?.code, wynik.blad?.message));
    return;
  }
  zaplecze.stan.wchlon(wynik.wynik.document);
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz('Dokument zapisany; wersja założona w repozytorium sesji.', true);
}

/**
 * „Wstaw do dokumentu" z paska narzędzi promptu — wynik operacji w miejsce kursora.
 *
 * Wstawienie nie jest przyjęciem propozycji: treść trafia do bufora edytora,
 * a propozycja czeka dalej na decyzję. Zrównanie obu czynności odebrałoby
 * Operatorowi możliwość wstawienia fragmentu i odrzucenia reszty.
 */
export function wstawWMiejsceKursora(
  stan: StanStudio,
  kontrolka: HTMLTextAreaElement,
  odpowiedz: WierszOdpowiedzi,
): void {
  const propozycja = stan.propozycja();
  if (propozycja === null || propozycja.tresc === '') {
    odpowiedz.pokaz('Nie ma wyniku operacji do wstawienia — najpierw uruchom operację w Tools Panel.', false);
    return;
  }
  const przed = kontrolka.value.slice(0, kontrolka.selectionStart);
  const po = kontrolka.value.slice(kontrolka.selectionEnd);
  stan.ustawTresc(`${przed}${propozycja.tresc}${po}`);
  odpowiedz.pokaz('Wynik operacji wstawiony w miejsce kursora; zapis zakłada wersję.', true);
}

/** Zapamiętuje zaznaczenie edytora; zakres pusty znaczy „cały dokument". */
export function zapamietajZaznaczenie(stan: StanStudio, kontrolka: HTMLTextAreaElement): void {
  const poczatek = kontrolka.selectionStart;
  const koniec = kontrolka.selectionEnd;
  stan.ustawZaznaczenie(poczatek === koniec ? null : { poczatek, koniec });
}
