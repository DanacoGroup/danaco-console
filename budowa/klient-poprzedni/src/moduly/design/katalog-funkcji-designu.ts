import { Command } from '../../../../shared/contract';
import type { WarstwaWidocznosci } from './warstwy-designu';

// Katalog funkcji modułu Design: pełny wykaz pozycji ze sposobem wykonania każdej w tej budowie.

/** Nazwy okien modułu Design, każde w jednym ustalonym brzmieniu, stosowanym w całym module bez odmian. */
export const OKNA = {
  chat: 'Chat Window',
  petla: 'Execution Loop Window',
  plansza: 'Design Board',
  zasoby: 'Assets Panel',
  kreator: 'Prompt Builder',
  podglad: 'Preview Window',
  tokeny: 'Tokens & System Panel',
} as const;

/**
 * Czym pozycja jest wykonywana.
 *
 * `komenda` — wykonuje ją komenda kontraktu wymieniona w polu `komendy`;
 * `okno`    — wykonuje ją samo okno, bez wywołania rdzenia;
 * `bez-drogi` — nie wykonuje jej ani okno, ani kontrakt.
 */
export type WykonanieFunkcji = 'komenda' | 'okno' | 'bez-drogi';

/** Jedna pozycja katalogu funkcji modułu Design wraz z grupą, opisem, oknem, warstwą widoczności i sposobem wykonania. */
export interface PozycjaKatalogu {
  /** Nazwa własna pozycji, zapisana dosłownie i niepodlegająca tłumaczeniu ani parafrazie. */
  readonly nazwa: string;
  /** Grupa tematyczna pozycji katalogu. */
  readonly grupa: string;
  /** Co pozycja robi. */
  readonly opis: string;
  /** Okno modułu, w którym pozycja jest osiągalna. */
  readonly okno: string;
  /** Warstwa widoczności pozycji. */
  readonly warstwa: WarstwaWidocznosci;
  /** Czym pozycja jest wykonywana. */
  readonly wykonanie: WykonanieFunkcji;
  /** Komendy kontraktu wykonujące pozycję; wypełnione wyłącznie przy `komenda`. */
  readonly komendy: readonly string[];
  /** Czego brakuje albo co robi okno zamiast tego; obowiązkowe przy `bez-drogi`. */
  readonly uwaga?: string;
}

const GENEROWANIE = 'Generowanie grafiki przez AI';
const RASTER = 'Edycja rastrowa';
const WEKTOR = 'Grafika i edycja wektorowa';
const MAKIETY = 'Projektowanie UI i makiety';
const TOKENY = 'Tokeny projektowe i system projektowy';
const KOLOR = 'Kolor';
const IKONY = 'Ikony i typografia';
const WYDANIE = 'Eksport, podgląd, handoff';
const MARKETING = 'Marketing, szablony, zasoby zewnętrzne';
const WSPOLPRACA = 'Współpraca, wersjonowanie, organizacja';

export const KATALOG_FUNKCJI: readonly PozycjaKatalogu[] = [
  // ── Generowanie grafiki przez AI ───────────────────────────────────────────
  {
    nazwa: 'Text-to-image',
    grupa: GENEROWANIE,
    opis: 'Generuje obraz z opisu skomponowanego w Prompt Builderze.',
    okno: OKNA.kreator,
    warstwa: 1,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetGenerate],
  },
  {
    nazwa: 'Image-to-image',
    grupa: GENEROWANIE,
    opis: 'Przekształca obraz referencyjny wedle promptu, z regulacją siły wpływu.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetGenerate],
    uwaga:
      'Referencję niesie pole referenceAssetId, siłę wpływu — pole creativity promptu. ' +
      'Osobnego parametru siły wobec referencji kontrakt nie ma.',
  },
  {
    nazwa: 'Inpainting (domalowanie fragmentu)',
    grupa: GENEROWANIE,
    opis: 'Zamienia zaznaczony obszar obrazu nowym, wpasowanym w otoczenie.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Żądanie generowania nie niesie ani maski, ani prostokąta obszaru, więc nie ma czym ' +
      'wskazać fragmentu do domalowania.',
  },
  {
    nazwa: 'Outpainting (rozszerzenie kadru)',
    grupa: GENEROWANIE,
    opis: 'Domalowuje treść poza pierwotną ramką obrazu.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Wymaga wskazania kierunku i wielkości dołożonego płótna wraz z zasobem bazowym; ' +
      'żądanie generowania nie ma na to ani jednego pola.',
  },
  {
    nazwa: 'Wariacje',
    grupa: GENEROWANIE,
    opis: 'Tworzy wiele wariantów tego samego zasobu do wyboru.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetGenerate],
    uwaga: 'Liczbę wariantów niesie pole variants promptu, powtarzalność — pole seed.',
  },
  {
    nazwa: 'Upscaling neuronowy',
    grupa: GENEROWANIE,
    opis: 'Powiększa obraz z odtworzeniem szczegółu.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageUpscale],
    uwaga: 'Kontrolki w oknie nie ma — czynność zleca się poleceniem w Chat Window.',
  },
  {
    nazwa: 'Usuwanie tła',
    grupa: GENEROWANIE,
    opis: 'Odcina obiekt od tła i tworzy kanał przezroczystości.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageBackgroundRemove],
    uwaga: 'Kontrolki w oknie nie ma — czynność zleca się poleceniem w Chat Window.',
  },
  {
    nazwa: 'Szablony promptu i historia poleceń',
    grupa: GENEROWANIE,
    opis:
      'Utrwala prompt strukturalny jako szablon do wielokrotnego użycia i oddaje wykaz ' +
      'poleceń wydanych w oknie wraz z zasobami, które z nich powstały.',
    okno: OKNA.kreator,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [
      Command.DesignPromptTemplateSave,
      Command.DesignPromptTemplateList,
      Command.DesignPromptHistoryList,
    ],
    uwaga:
      'Prompt bez ani jednego zasobu zostaje w historii: kanał bywa odmawiał, a wtedy zapis ' +
      'jest zapisem próby — po to sięga się do historii.',
  },
  {
    nazwa: 'Transfer stylu',
    grupa: GENEROWANIE,
    opis: 'Przenosi styl obrazu-wzorca na treść obrazu-bazy.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Żądanie zna jeden zasób referencyjny i nie rozróżnia wzorca stylu od obrazu bazowego, ' +
      'a transfer wymaga obu naraz.',
  },
  {
    nazwa: 'Wektoryzacja (raster→SVG)',
    grupa: GENEROWANIE,
    opis: 'Zamienia bitmapę na ścieżki wektorowe.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageVectorize],
    uwaga:
      'Komenda przyjmuje sposób obrysu, liczbę barw i siłę upraszczania ścieżek. Kontrolki ' +
      'w oknie nie ma — czynność zleca się poleceniem w Chat Window.',
  },
  {
    nazwa: 'Generowanie ikon i logotypów',
    grupa: GENEROWANIE,
    opis: 'Tworzy spójny zestaw ikon i znaków w jednym stylu.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetGenerate],
    uwaga:
      'Jedno zlecenie niesie jeden temat, więc zestaw powstaje z wielu zleceń; wymuszenia ' +
      'siatki i grubości obrysu żądanie nie ma czym wyrazić.',
  },
  {
    nazwa: 'Rozdzielenie na warstwy',
    grupa: GENEROWANIE,
    opis: 'Rozkłada wygenerowany obraz na obiekty i warstwy edytowalne.',
    okno: OKNA.plansza,
    warstwa: 4,
    wykonanie: 'komenda',
    komendy: [Command.ImageLayersSplit],
    uwaga:
      'Komenda oddaje po jednym zasobie na rozpoznany obiekt. Bez zainstalowanego silnika ' +
      'segmentacji odmawia, nazywając brak — nigdy nie oddaje całego obrazu jako jednej ' +
      'warstwy, udając rozkład.',
  },
  {
    nazwa: 'Naprawa twarzy i detalu',
    grupa: GENEROWANIE,
    opis: 'Poprawia zniekształcone twarze oraz drobne detale generacji.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageUpscale],
    uwaga: 'Osobnym przebiegiem powiększenia — pole faces żądania. Kontrolki w oknie nie ma.',
  },
  {
    nazwa: 'Sterowanie kompozycją (mapy sterujące)',
    grupa: GENEROWANIE,
    opis: 'Wymusza pozę, szkielet, głębię albo krawędzie w generacji wedle szkicu.',
    okno: OKNA.kreator,
    warstwa: 4,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Mapa sterująca jest osobnym obrazem wejściowym, a żądanie generowania zna jeden.',
  },

  // ── Edycja rastrowa ────────────────────────────────────────────────────────
  {
    nazwa: 'Kadrowanie i prostowanie',
    grupa: RASTER,
    opis: 'Przycina, obraca, prostuje horyzont, zmienia proporcje.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageTransform],
    uwaga: 'Kontrolki w oknie nie ma — czynność zleca się poleceniem w Chat Window.',
  },
  {
    nazwa: 'Skalowanie i resampling',
    grupa: RASTER,
    opis: 'Zmienia rozmiar z zachowaniem albo złamaniem proporcji.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageTransform],
    uwaga: 'Wyboru algorytmu resamplingu żądanie nie niesie — rozstrzyga go rdzeń.',
  },
  {
    nazwa: 'Warstwy rastrowe',
    grupa: RASTER,
    opis: 'Niezależne warstwy z trybami mieszania i krycia.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.ImageCompose],
    uwaga:
      'Składanie obrazu na obrazie wraz z trybem mieszania i kryciem wykonuje komenda, a jej ' +
      'wynikiem jest NOWY zasób. Warstwa kompozycji tego nie wyraża: niesie położenie, rozmiar, ' +
      'kolejność, blokadę i adnotację — ani trybu mieszania, ani krycia. Mieszanie jest więc ' +
      'czynnością na zasobach, a nie właściwością kanwy.',
  },
  {
    nazwa: 'Maski warstw',
    grupa: RASTER,
    opis: 'Nieniszcząca maska przezroczystości, malowana albo wzięta z zaznaczenia.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Maska jest osobnym obrazem wiązanym z warstwą; kontrakt nie zna takiego wiązania.',
  },
  {
    nazwa: 'Korekcja barwna',
    grupa: RASTER,
    opis: 'Jasność, kontrast, nasycenie i wyrównanie poziomów.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageAdjust],
    uwaga: 'Temperatury barwowej, krzywych i osobnego sterowania barwą kontrakt nie zna.',
  },
  {
    nazwa: 'Retusz (klonowanie i leczenie)',
    grupa: RASTER,
    opis: 'Klonuje i wygładza fragmenty, usuwa niedoskonałości.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Poprawka obrazu obejmuje cały obraz jednym natężeniem; pędzla ani obszaru nie ma czym ' +
      'wskazać.',
  },
  {
    nazwa: 'Filtry i efekty',
    grupa: RASTER,
    opis: 'Rozmycie, wyostrzenie, ziarno, winieta, cień, poświata.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageAdjust],
    uwaga:
      'Rozmycie, wyostrzenie, odszumienie i skala szarości są w wyliczeniu poprawek; ziarna, ' +
      'winiety, cienia i poświaty w nim nie ma.',
  },
  {
    nazwa: 'Zaznaczanie obiektu',
    grupa: RASTER,
    opis: 'Automatyczne zaznaczenie obiektu jednym kliknięciem.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.ImageLayersSplit],
    uwaga:
      'Segmentacja oddaje obiekty jako osobne zasoby, więc „zaznaczenie" jest tu rozkładem ' +
      'obrazu, a nie zaznaczeniem w kanwie. Kanwa nadal zaznacza warstwy, nie obiekty wewnątrz ' +
      'obrazu — jednym kliknięciem to nie jest.',
  },
  {
    nazwa: 'Korekcja perspektywy',
    grupa: RASTER,
    opis: 'Prostuje zniekształcenia perspektywiczne i obiektywu.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Wyliczenie przekształceń geometrycznych obejmuje skalowanie, kadr, obrót, odbicia ' +
      'i miniaturę; przekształcenia homograficznego w nim nie ma.',
  },
  {
    nazwa: 'Tryb wsadowy edycji',
    grupa: RASTER,
    opis: 'Stosuje ten sam zestaw operacji do wielu zasobów naraz.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Jedno żądanie obrazu dotyczy jednego zasobu, a kolejki modułu Automations kontrakt ' +
      'z tym modułem nie wiąże; zestawu operacji nie ma też czym zapisać jako presetu.',
  },

  // ── Grafika i edycja wektorowa ─────────────────────────────────────────────
  {
    nazwa: 'Narzędzie pióra (ścieżki Béziera)',
    grupa: WEKTOR,
    opis: 'Rysuje i edytuje krzywe oraz węzły ścieżek.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Kanwa układa warstwy, nie węzły; kompozycja nie ma bytu ścieżki, więc rysunku nie ma ' +
      'gdzie zapisać.',
  },
  {
    nazwa: 'Kształty podstawowe i złożone',
    grupa: WEKTOR,
    opis: 'Prostokąty, elipsy, wielokąty, gwiazdy, zaokrąglenia.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Elementy pomocnicze kanwy są prostokątnymi warstwami bez zasobu; kształtu innego niż ' +
      'prostokąt warstwa nie wyraża.',
  },
  {
    nazwa: 'Operacje logiczne (boolean)',
    grupa: WEKTOR,
    opis: 'Suma, różnica, przecięcie i wykluczenie ścieżek.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga bytu ścieżki, którego kompozycja nie zna.',
  },
  {
    nazwa: 'Obrys i wypełnienie',
    grupa: WEKTOR,
    opis: 'Grubość, zakończenia, gradient, deseń, reguła wypełnienia.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Warstwa kompozycji nie niesie ani obrysu, ani wypełnienia.',
  },
  {
    nazwa: 'Tekst na ścieżce i obrys tekstu',
    grupa: WEKTOR,
    opis: 'Układa tekst wzdłuż krzywej i zamienia go w kontury.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga bytu ścieżki i bytu tekstu; kompozycja nie ma żadnego z nich.',
  },
  {
    nazwa: 'Optymalizacja SVG',
    grupa: WEKTOR,
    opis: 'Czyści i minimalizuje kod SVG bez utraty jakości.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Konwersja obrazu dotyczy formatów rastrowych; zasobu wektorowego nie ma czym przetworzyć ' +
      'ani odczytać po treści.',
  },
  {
    nazwa: 'Symbole i instancje',
    grupa: WEKTOR,
    opis: 'Definicja elementu wielokrotnego z propagacją zmian.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Kompozycja nie zna wiązania warstwy z definicją wspólną.',
  },
  {
    nazwa: 'Siatka i przyciąganie precyzyjne',
    grupa: WEKTOR,
    opis: 'Rozmieszczanie do siatki, prowadnic i punktów.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Siatka pomocnicza i wyrównanie warstw dzieją się na kanwie; do rdzenia jedzie dopiero ' +
      'gotowy układ. Przyciągania do prowadnic kanwa nie prowadzi.',
  },
  {
    nazwa: 'Eksport wektorowy',
    grupa: WEKTOR,
    opis: 'Wydaje czysty SVG oraz PDF wektorowy.',
    okno: OKNA.podglad,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetExport],
    uwaga:
      'Eksport przyjmuje format wektorowy i dokumentowy obok rastrowych. Kontrolki w oknie nie ' +
      'ma — czynność zleca się poleceniem w Chat Window.',
  },

  // ── Projektowanie UI i makiety ─────────────────────────────────────────────
  {
    nazwa: 'Ramki i obszary robocze',
    grupa: MAKIETY,
    opis: 'Wydzielone ekrany o zdefiniowanych rozmiarach urządzeń.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'okno',
    komendy: [Command.DesignBoardUpdate],
    uwaga:
      'Ramkę urządzenia zakłada przybornik kanwy, a jej wymiary jadą do rdzenia polami ' +
      'width i height warstwy.',
  },
  {
    nazwa: 'Komponenty i warianty',
    grupa: MAKIETY,
    opis: 'Element wielokrotny z wariantami stanu.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Rejestru komponentów i właściwości wariantu kontrakt nie zna.',
  },
  {
    nazwa: 'Auto-layout',
    grupa: MAKIETY,
    opis: 'Automatyczne rozmieszczanie z odstępami i dopasowaniem do treści.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Kanwa rozmieszcza warstwy równomiernie i układa je szablonem; odstępu i wyściółki jako ' +
      'trwałych właściwości warstwa nie niesie, więc układ jest jednorazowy.',
  },
  {
    nazwa: 'Więzy responsywne',
    grupa: MAKIETY,
    opis: 'Zachowanie elementów przy zmianie rozmiaru ramki.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Reguły kotwiczenia nie ma czym zapisać przy warstwie.',
  },
  {
    nazwa: 'Prototypowanie i przejścia',
    grupa: MAKIETY,
    opis: 'Łączy ekrany interakcjami i animacjami przejść.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Grafu połączeń ekranów kompozycja nie niesie.',
  },
  {
    nazwa: 'Wireframe niskiej wierności',
    grupa: MAKIETY,
    opis: 'Szybkie makiety szkicowe z biblioteką prostych elementów.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [Command.DesignBoardUpdate],
    uwaga: 'Elementy pomocnicze kanwy są warstwami bez zasobu i jadą do rdzenia z układem.',
  },
  {
    nazwa: 'Biblioteka UI',
    grupa: MAKIETY,
    opis: 'Gotowe komponenty do składania makiet.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga rejestru komponentów, którego kontrakt nie zna.',
  },
  {
    nazwa: 'Generowanie makiety z opisu',
    grupa: MAKIETY,
    opis: 'Tworzy szkielet ekranu z polecenia tekstowego.',
    okno: OKNA.kreator,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetGenerate],
    uwaga:
      'Powstaje OBRAZ makiety, nie drzewo komponentów — kanał obrazowy oddaje bajty obrazu ' +
      'i niczego innego.',
  },
  {
    nazwa: 'Import ze zrzutu ekranu',
    grupa: MAKIETY,
    opis: 'Odtwarza edytowalną makietę z obrazu istniejącego interfejsu.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Wgranie wnosi obraz jako zasób i na tym się kończy; rozpoznania układu kontrakt nie ' +
      'oddaje.',
  },
  {
    nazwa: 'Siatki układu',
    grupa: MAKIETY,
    opis: 'Kolumny, rynny, moduły, siatka bazowa.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Kanwa prowadzi siatkę bazową o stałym module; definicji kolumn i rynien nie ma czym ' +
      'ani ustawić, ani zapisać.',
  },

  // ── Tokeny projektowe i system projektowy ──────────────────────────────────
  {
    nazwa: 'Tokeny kolorów',
    grupa: TOKENY,
    opis: 'Definicja nazwanych barw semantycznych.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignTokensetSave, Command.DesignTokensetList],
    uwaga:
      'Drzewo czyta wartości z motywu obowiązującego, a zestaw własny ma już gdzie zamieszkać: ' +
      'zestaw żetonów jest osobnym bytem wraz ze wskazaniem motywu, którego dotyczy. Okno tej ' +
      'drogi jeszcze nie wywołuje.',
  },
  {
    nazwa: 'Tokeny typografii',
    grupa: TOKENY,
    opis: 'Skala rozmiarów, krój, interlinia, grubość, odstęp liter.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignTokensetSave, Command.DesignTokensetList],
    uwaga: 'Jak przy barwach: odczyt z motywu w oknie, trwały zapis komendą zestawu żetonów.',
  },
  {
    nazwa: 'Tokeny odstępów i promieni',
    grupa: TOKENY,
    opis: 'Skala odstępów, promieni, wymiarów i warstw.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignTokensetSave, Command.DesignTokensetList],
    uwaga: 'Jak przy barwach: odczyt z motywu w oknie, trwały zapis komendą zestawu żetonów.',
  },
  {
    nazwa: 'Motyw jasny i ciemny',
    grupa: TOKENY,
    opis: 'Warianty żetonów dla obu motywów wraz z przełącznikiem podglądu.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Panel przełącza motyw podglądu i pokazuje obie kolumny wartości. Wariantów marki ' +
      'poza dwoma motywami produktu nie ma.',
  },
  {
    nazwa: 'Eksport tokenów do kodu',
    grupa: TOKENY,
    opis: 'Wydaje żetony jako zmienne CSS, SCSS, konfigurację Tailwind i moduł JavaScript.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignTokensetExport],
    uwaga:
      'Cztery postacie składa i oddaje przeglądarka bez rdzenia. Komenda dokłada dwie, których ' +
      'przeglądarka złożyć nie może — dla systemów mobilnych — oraz wydanie idące DO MODUŁU ' +
      'zamiast do pliku. Okno wywołuje dziś tylko drogę własną.',
  },
  {
    nazwa: 'Import tokenów',
    grupa: TOKENY,
    opis: 'Wczytuje istniejące żetony do modułu.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignTokensetImport],
    uwaga:
      'Okno zestawia wczytany zapis z żetonami motywu i pokazuje różnicę. Komenda dokłada ' +
      'utrwalenie go jako zestawu modułu wraz z wykazem ról spoza systemu produktu. Nadpisania ' +
      'motywu nie robi ani okno, ani rdzeń: motyw jest własnością powłoki.',
  },
  {
    nazwa: 'Powiązanie tokenów z komponentami',
    grupa: TOKENY,
    opis: 'Wskazuje komponenty korzystające z danego żetonu.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga: 'Panel przeszukuje arkusze produktu i wypisuje selektory, w których żeton występuje.',
  },
  {
    nazwa: 'Dokumentacja systemu',
    grupa: TOKENY,
    opis: 'Generuje przewodnik: kolory, typografia, komponenty.',
    okno: OKNA.tokeny,
    warstwa: 4,
    wykonanie: 'komenda',
    komendy: [Command.DesignStyleguidePublish],
    uwaga:
      'Przewodnik składa przeglądarka i oddaje go plikiem. Wydanie go do modułu docelowego ' +
      'wykonuje komenda, zakładając zasób przewodnika w magazynie — okno tej drogi jeszcze ' +
      'nie wywołuje.',
  },

  // ── Kolor ─────────────────────────────────────────────────────────────────
  {
    nazwa: 'Generator palet',
    grupa: KOLOR,
    opis: 'Tworzy paletę z reguły harmonii.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Rachunek harmonii jest wykonalny w przeglądarce, ale w tej budowie go nie ma; panel ' +
      'prowadzi kontrast i symulację widzenia, generatora palet nie.',
  },
  {
    nazwa: 'Ekstrakcja palety z obrazu',
    grupa: KOLOR,
    opis: 'Wyciąga dominujące barwy z zasobu.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Bajty obrazu mają już drogę — komenda oddająca treść zasobu jest w kontrakcie. Samego ' +
      'wyciągania barw dominujących nie robi ani ona, ani okno, więc pozycja czeka na rachunek, ' +
      'nie na drogę.',
  },
  {
    nazwa: 'Kontroler kontrastu WCAG',
    grupa: KOLOR,
    opis: 'Liczy współczynnik kontrastu i ocenę dla par barw.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Pary i progi pochodzą z wykazu progów kontrastu produktu; rachunek luminancji wykonuje ' +
      'przeglądarka na wartościach żetonów obowiązującego motywu.',
  },
  {
    nazwa: 'Symulacja wad wzroku',
    grupa: KOLOR,
    opis: 'Podgląd barw w protanopii, deuteranopii i tritanopii.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga: 'Symulacja obejmuje próbki żetonów; obrazów nie obejmuje, bo nie ma po nie drogi.',
  },
  {
    nazwa: 'Edytor gradientów',
    grupa: KOLOR,
    opis: 'Gradienty liniowe, promieniste i kątowe z wieloma stopniami.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Gradient nie jest żetonem, który panel czyta, ani polem warstwy kompozycji — nie ma go ' +
      'gdzie zapisać.',
  },
  {
    nazwa: 'Konwersja przestrzeni barw',
    grupa: KOLOR,
    opis: 'Zapis barwy w kilku notacjach obok siebie.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Wiersz żetonu barwnego podaje zapis szesnastkowy wraz ze składowymi i luminancją. ' +
      'Przestrzeni Lab i CMYK panel nie liczy.',
  },
  {
    nazwa: 'Sprawdzian dostępności kolorem',
    grupa: KOLOR,
    opis: 'Wskazuje pary żetonów łamiące próg kontrastu w całym systemie.',
    okno: OKNA.tokeny,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga: 'Wynikiem jest wykaz par poniżej progu wraz z ich zmierzonym współczynnikiem.',
  },

  // ── Ikony i typografia ─────────────────────────────────────────────────────
  {
    nazwa: 'Biblioteka ikon',
    grupa: IKONY,
    opis: 'Przeszukiwalny katalog ikon do wstawienia.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Zestaw ikon produktu należy do powłoki, a wstawienia ikony na kanwę jako warstwy ' +
      'kompozycja nie wyraża — warstwa niesie zasób albo nic.',
  },
  {
    nazwa: 'Edytor ikony na siatce',
    grupa: IKONY,
    opis: 'Rysuje i poprawia ikonę na siatce z wyrównaniem do pikseli.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga edytora ścieżek, którego kanwa nie prowadzi.',
  },
  {
    nazwa: 'Generowanie zestawu ikon',
    grupa: IKONY,
    opis: 'Tworzy spójny komplet ikon w jednym stylu z listy pojęć.',
    okno: OKNA.kreator,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetGenerate],
    uwaga: 'Jedno zlecenie niesie jedno pojęcie; listy pojęć żądanie nie przyjmuje.',
  },
  {
    nazwa: 'Font ikon i sprite',
    grupa: IKONY,
    opis: 'Pakuje ikony w font ikon albo sprite z symbolami.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga wielu zasobów wejściowych i jednego pliku wyjściowego — kontrakt nie ma obu.',
  },
  {
    nazwa: 'Generator faviconów',
    grupa: IKONY,
    opis: 'Wydaje komplet faviconów i ikon aplikacji ze źródła.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetExportBatch],
    uwaga:
      'Eksport zbiorczy przyjmuje komplet skal i format ikony systemowej w jednym wywołaniu. ' +
      'Manifestu aplikacji nie składa — to jest plik opisowy, nie obraz, i pozostaje bez drogi.',
  },
  {
    nazwa: 'Dobór par krojów',
    grupa: IKONY,
    opis: 'Proponuje pasujące zestawienia krojów nagłówka i tekstu.',
    okno: OKNA.tokeny,
    warstwa: 4,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Panel czyta trzy kroje żetonów typografii i nic poza nimi; katalogu krojów kontrakt ' +
      'nie oddaje.',
  },
  {
    nazwa: 'Podgląd i osadzenie webfontów',
    grupa: IKONY,
    opis: 'Podgląd kroju w tekście próbnym i wydanie reguły osadzenia.',
    okno: OKNA.tokeny,
    warstwa: 4,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga wniesienia pliku kroju; wgranie zasobu przyjmuje obrazy, nie kroje pisma.',
  },
  {
    nazwa: 'Inspektor glifów',
    grupa: IKONY,
    opis: 'Przegląd znaków, ligatur i wariantów typograficznych kroju.',
    okno: OKNA.tokeny,
    warstwa: 4,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga: 'Wymaga odczytu pliku kroju, którego moduł nie ma skąd wziąć.',
  },

  // ── Eksport, podgląd, handoff ──────────────────────────────────────────────
  {
    nazwa: 'Odczyt treści zasobu z magazynu',
    grupa: WYDANIE,
    opis:
      'Oddaje bajty zasobu wraz z ich wielkością i sumą kontrolną — jedyna droga do treści, ' +
      'bo pole uri zasobu jest ścieżką w systemie plików rdzenia.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetContentGet],
    uwaga:
      'Rdzeń odmawia zasobu nieznanego i zasobu bez treści pod sumą kontrolną; bajtów ' +
      'zastępczych nie oddaje żadną drogą.',
  },
  {
    nazwa: 'Eksport wieloformatowy',
    grupa: WYDANIE,
    opis: 'Wydaje zasób w wybranym formacie.',
    okno: OKNA.podglad,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetExport, Command.ImageConvert],
    uwaga:
      'Dwie różne czynności, nie jedna: eksport oddaje Operatorowi bajty gotowe do zapisania ' +
      'poza produktem, a konwersja zakłada NOWY zasób w magazynie i tam się kończy. Opracowanie ' +
      'ma na myśli pierwszą.',
  },
  {
    nazwa: 'Kompresja i optymalizacja',
    grupa: WYDANIE,
    opis: 'Redukuje wagę pliku sterując jakością kompresji.',
    okno: OKNA.podglad,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.ImageConvert],
    uwaga: 'Odpowiedź podaje rozmiar po konwersji i oszczędność wobec źródła.',
  },
  {
    nazwa: 'Skalowanie @1x/@2x/@3x',
    grupa: WYDANIE,
    opis: 'Generuje warianty gęstości pikseli jednym poleceniem.',
    okno: OKNA.podglad,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetExportBatch],
    uwaga:
      'Komplet skal idzie jednym wywołaniem, a nazwa wariantu gęstości powstaje z krotnosci — ' +
      '„jednym poleceniem" jest już prawdą.',
  },
  {
    nazwa: 'Cięcie na fragmenty',
    grupa: WYDANIE,
    opis: 'Eksportuje wskazane obszary makiety jako osobne pliki.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageTransform],
    uwaga:
      'Kadr wycina obszar z JEDNEGO obrazu; ramek eksportu opisanych na kompozycji kontrakt ' +
      'nie zna, a wynik zostaje w magazynie zamiast wyjść plikiem.',
  },
  {
    nazwa: 'Handoff i inspekcja',
    grupa: WYDANIE,
    opis: 'Udostępnia wymiary, odstępy i zapis stylu elementu.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Inspektor właściwości podaje położenie i wymiary warstwy zaznaczonej wraz z gotowym ' +
      'zapisem reguł stylu. Barw i typografii warstwa nie niesie, więc ich w zapisie nie ma.',
  },
  {
    nazwa: 'Eksport kodu widoku',
    grupa: WYDANIE,
    opis: 'Zamienia makietę albo element w kod widoku.',
    okno: OKNA.plansza,
    warstwa: 4,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Wymaga drzewa komponentów, którego kompozycja nie niesie, i drogi wydania do modułu ' +
      'Apps, której kontrakt nie zna.',
  },
  {
    nazwa: 'Eksport tablicy jako PDF lub obraz',
    grupa: WYDANIE,
    opis: 'Zapisuje całą kompozycję do jednego pliku.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignBoardExport],
    uwaga:
      'Komenda wyrysowuje całą kompozycję albo wskazany obszar. Wyrys wymaga jednak TREŚCI ' +
      'zasobów leżących w warstwach, więc jest zależny od tej samej drogi po bajty, która ' +
      'gasi dziś podgląd.',
  },
  {
    nazwa: 'Osadzenie w ramce urządzenia',
    grupa: WYDANIE,
    opis: 'Wstawia zasób w ramkę telefonu, laptopa albo wizytówki.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageCompose],
    uwaga:
      'Złożenie dwóch obrazów wykonuje komenda: zasób wchodzi w ramkę wraz z położeniem, skalą ' +
      'i kryciem. Biblioteki ramek urządzeń moduł nie ma — ramkę trzeba wnieść jako zasób.',
  },
  {
    nazwa: 'Eksport animacji lekkiej',
    grupa: WYDANIE,
    opis: 'Wydaje animację jako plik animowany.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Rodzaje zasobu obejmują obraz, wektor, kompozycję, dokument, dźwięk, film i archiwum; ' +
      'animacji lekkiej wśród nich nie ma, a klatek nie ma z czego złożyć.',
  },
  {
    nazwa: 'Metadane i profil barwny',
    grupa: WYDANIE,
    opis: 'Odczytuje i osadza metadane oraz profil barw przy wydaniu.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageInspect],
    uwaga:
      'Odczyt jest — przestrzeń barw i metadane wracają w odpowiedzi. Osadzenia ani ' +
      'oczyszczenia metadanych kontrakt nie zna.',
  },

  // ── Marketing, szablony, zasoby zewnętrzne ─────────────────────────────────
  {
    nazwa: 'Szablony formatów społecznościowych',
    grupa: MARKETING,
    opis: 'Gotowe rozmiary postów, relacji i okładek.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga:
      'Przybornik kanwy układa warstwy w formatach społecznościowych. Rozmiary są układem ' +
      'kanwy, nie wydaniem pliku — eksportu nadal nie ma.',
  },
  {
    nazwa: 'Zestawy rozmiarów kampanii',
    grupa: MARKETING,
    opis: 'Generuje ten sam projekt w wielu formatach reklamowych naraz.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetExportBatch],
    uwaga:
      'Eksport zbiorczy wydaje ten sam zasób w komplecie skal jednym wywołaniem. Przeskalowania ' +
      'UKŁADU do innych proporcji to nie daje: okno prowadzi jedną kanwę, a zestawu kilku ramek ' +
      'wynikowych nie ma gdzie zapisać.',
  },
  {
    nazwa: 'Biblioteka szablonów brandingowych',
    grupa: MARKETING,
    opis: 'Arkusze brandingu, tablice nastroju i siatki porównawcze.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'okno',
    komendy: [],
    uwaga: 'Szablony układu stoją w przyborniku kanwy i działają na warstwach kompozycji.',
  },
  {
    nazwa: 'Import zasobów zewnętrznych',
    grupa: MARKETING,
    opis: 'Wyszukuje i wstawia zasoby z bibliotek zewnętrznych.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Wgranie wnosi plik wskazany przez Operatora; odpytania biblioteki zewnętrznej kontrakt ' +
      'nie zna.',
  },
  {
    nazwa: 'Placeholdery treści',
    grupa: MARKETING,
    opis: 'Wstawia treści próbne zasilające makiety.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Kontrakt nie zna generatora treści próbnych, a kierunek projektowy produktu zabrania ' +
      'zmyślonych osób i metryk — treść próbna musiałaby pochodzić z domeny i być oznaczona.',
  },
  {
    nazwa: 'Znak wodny i branding wsadowy',
    grupa: MARKETING,
    opis: 'Nakłada logo albo znak wodny na zestaw zasobów.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.ImageCompose],
    uwaga:
      'Nałożenie znaku na jeden zasób wykonuje komenda składania. WSADU to nie daje: jedno ' +
      'wywołanie dotyczy jednej pary obrazów, więc zestaw to tyle wywołań, ile zasobów.',
  },

  // ── Współpraca, wersjonowanie, organizacja ─────────────────────────────────
  {
    nazwa: 'Warstwy i drzewo obiektów',
    grupa: WSPOLPRACA,
    opis: 'Panel warstw z widocznością, blokadą i zagnieżdżeniem.',
    okno: OKNA.plansza,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [Command.DesignBoardUpdate, Command.DesignBoardList],
    uwaga:
      'Kolejność, blokada i adnotacja jadą do rdzenia. Zagnieżdżenia w grupy warstwa nie ' +
      'wyraża — nie ma pola wskazującego warstwę nadrzędną.',
  },
  {
    nazwa: 'Wersjonowanie kompozycji',
    grupa: WSPOLPRACA,
    opis: 'Historia stanów tablicy, nazwane wersje, powrót, porównanie.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [
      Command.DesignBoardVersionSave,
      Command.DesignBoardVersionList,
      Command.DesignBoardVersionRestore,
    ],
    uwaga:
      'Wersja jest już osobnym bytem: zapis utrwala układ pod nazwą, wykaz oddaje ciąg postaci, ' +
      'a przywrócenie zakłada wersję z układu sprzed cofnięcia, żeby porzucony stan nie przepadł. ' +
      'Porównania dwóch wersji wprost kontrakt nie ma — zestawia się je odczytem obu.',
  },
  {
    nazwa: 'Diff wizualny zasobów',
    grupa: WSPOLPRACA,
    opis: 'Nakłada dwie wersje i podświetla różnice.',
    okno: OKNA.podglad,
    warstwa: 3,
    wykonanie: 'bez-drogi',
    komendy: [],
    uwaga:
      'Porównanie wariantów zestawia dziś pola opisowe. Bajty obu stron mają już drogę — komenda ' +
      'oddająca treść zasobu jest w kontrakcie — więc brakuje samego nałożenia i wskazania ' +
      'różnic, a nie dostępu do obrazów.',
  },
  {
    nazwa: 'Komentarze i adnotacje',
    grupa: WSPOLPRACA,
    opis: 'Przypięte uwagi do elementów, wątki, oznaczenia osób.',
    okno: OKNA.plansza,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignAnnotationSet, Command.DesignAnnotationList, Command.DesignBoardUpdate],
    uwaga:
      'Adnotacja jest już osobnym bytem z autorem, wskazaniem nadrzędnym — tak powstaje wątek — ' +
      'i znacznikiem zamknięcia. Pole adnotacji warstwy zostaje obok, jako jedno zdanie jadące ' +
      'razem z układem; oznaczeń osób kontrakt nadal nie zna.',
  },
  {
    nazwa: 'Kursory i obecność',
    grupa: WSPOLPRACA,
    opis: 'Widoczność kursorów współpracujących osób na tablicy.',
    okno: OKNA.plansza,
    warstwa: 1,
    wykonanie: 'komenda',
    komendy: [Command.DesignPresenceReport],
    uwaga:
      'Zgłoszenie obecności jest ULOTNE — nie zapisuje się w bazie, bo położenie kursora sprzed ' +
      'godziny nie jest wiedzą o niczym. Rdzeń rozgłasza je pozostałym osobnym zdarzeniem ' +
      'obecności, więc kursory nie jadą przez zdarzenie zasobu.',
  },
  {
    nazwa: 'Tagi, kolekcje, wyszukiwanie',
    grupa: WSPOLPRACA,
    opis: 'Etykietowanie i przeszukiwanie zasobów po nazwie, etykiecie i prompcie.',
    okno: OKNA.zasoby,
    warstwa: 2,
    wykonanie: 'komenda',
    komendy: [
      Command.DesignAssetTagSet,
      Command.DesignAssetList,
      Command.DesignCollectionCreate,
      Command.DesignCollectionAssign,
      Command.DesignCollectionList,
    ],
    uwaga:
      'Etykiety i kolekcje są już dwoma osobnymi bytami: etykieta jest słowem, kolekcja ma nazwę, ' +
      'opis i porządek, a przypisanie jest dokładką albo odjęciem, nie zastąpieniem. Wyszukiwanie ' +
      'po treści promptu nadal nie jest polem żądania wykazu — fraza zawęża wyłącznie wykaz już ' +
      'wczytany, a pole frazy zgłoszono jako rozszerzenie, które do kontraktu nie weszło.',
  },
  {
    nazwa: 'Metadane pochodzenia (prowenancja)',
    grupa: WSPOLPRACA,
    opis: 'Zapisuje model, prompt, ziarno, datę i łańcuch edycji zasobu.',
    okno: OKNA.zasoby,
    warstwa: 3,
    wykonanie: 'komenda',
    komendy: [Command.DesignAssetList, Command.DesignPromptHistoryList],
    uwaga:
      'Wykaz promptów wnosi przekład, którego brakowało: prompt ma kod kontraktu i wymienia ' +
      'zasoby, które z niego powstały, wraz z kanałem i czasem wydania. Puste pole promptu przy ' +
      'zasobie przestaje więc być ślepą uliczką. Łańcucha edycji zasób nadal nie niesie — ' +
      'wskazuje wyłącznie pierwowzór wariantu.',
  },
];

/** Grupy katalogu funkcji w kolejności pierwszego wystąpienia pozycji, bez powtórzeń jednej grupy dwukrotnie. */
export function grupyKatalogu(): readonly string[] {
  return [...new Set(KATALOG_FUNKCJI.map((pozycja) => pozycja.grupa))];
}

/** Liczba pozycji katalogu o danym sposobie wykonania, liczona z wykazu, nie zapisana nigdzie na stałe jako osobna wartość. */
export function liczbaWedlugWykonania(wykonanie: WykonanieFunkcji): number {
  return KATALOG_FUNKCJI.filter((pozycja) => pozycja.wykonanie === wykonanie).length;
}
