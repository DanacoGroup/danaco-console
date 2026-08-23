import {
  utworzOknoWarsztatuDesignu,
  type NastawyWarsztatuDesignu,
  type OknoWarsztatuDesignu,
} from './okno-warsztatu-designu';
import type { Kanal } from '../../protokol/kanal';
import type { StanDesignu } from './stan-designu';
import { utworzZrodloWarsztatowDesignu } from './zrodlo-warsztatow-designu';

/**
 * Pięć warsztatów modułu Design — nastawy okien i ich złożenie.
 *
 * Bez tych okien siedemdziesiąt pięć komend byłoby funkcjami, których Operator
 * nie ma. Nastawy stoją danymi, nie pięcioma plikami po jednym oknie: różnią się
 * kodem katalogu rdzenia, tytułem, rolą, grupą czynności i zdaniem objaśnienia,
 * a poza tym są tym samym oknem (`okno-warsztatu-designu.ts`).
 *
 * Kody są kodami KATALOGU RDZENIA (migracja 340). Bez wiersza w katalogu klient
 * postawiłby okno, o którym rdzeń nie wie: `module.list` zaniżałby zakres modułu,
 * a pas uczciwości meldowałby „poza katalogiem".
 */
export const NASTAWY_WARSZTATOW_DESIGNU: readonly NastawyWarsztatuDesignu[] = [
  {
    kod: 'design.photo-workshop',
    tytul: 'Warsztat fotografii',
    rola: 'pomocnicze',
    grupa: 'fotografia',
    objasnienie:
      'Obróbka zdjęć od kadru do wsadu: kadrowanie i prostowanie, przekształcenia, ' +
      'rozdzielczość, powiększenie, poprawa jakości, korekcje barwne, filtry, retusz, ' +
      'domalowanie, rozszerzenie kadru, odcięcie tła, zaznaczenie, maski, warstwy, ' +
      'obrysowanie konturów, nastawy, wsad, metadane i łańcuch edycji. Podstawa każdej ' +
      'czynności liczy się rachunkiem WKOMPILOWANYM w rdzeń — bez ani jednego programu ' +
      'zewnętrznego, więc czynność działa u Operatora zawsze. Cztery czynności mają ' +
      'wariant lepszy od rachunku i odpowiedź mówi, którą drogą policzyła. Oryginał ' +
      'zostaje NIETKNIĘTY: każda obróbka zakłada wariant.',
  },
  {
    kod: 'design.vector-workshop',
    tytul: 'Warsztat wektora',
    rola: 'pomocnicze',
    grupa: 'wektor',
    objasnienie:
      'Pióro i węzły, kształty podstawowe, operacje logiczne, tekst na ścieżce i w ' +
      'konturach, czyszczenie zapisu, symbole oraz wydanie SVG, PDF i EPS. Kształt ' +
      'powstaje OD RAZU jako węzły ścieżki, więc da się go ciągnąć piórem od pierwszej ' +
      'chwili. Zamiana tekstu w kontury jest nieodwracalna dla wyniku, dlatego tekst ' +
      'źródłowy zostaje w nazwie ścieżki.',
  },
  {
    kod: 'design.print-workshop',
    tytul: 'Warsztat druku',
    rola: 'pomocnicze',
    grupa: 'druk',
    objasnienie:
      'Profile wydania, nośniki, kontrola przeddrukowa, wydanie do druku, kafle wielkiego ' +
      'formatu, wykresy i schematy. Kontrola oddaje zastrzeżenia w kolejności WAGI, a ' +
      'wada o wadze błędu ODMAWIA wydania: plik nie do druku wydany jako gotowy do druku ' +
      'jedzie do drukarni i kosztuje nakład. Pominięcie kontroli jest jawnym wyborem ' +
      'Operatora i wraca w odpowiedzi. Wykres bierze SERIE DANYCH, nie obraz, więc da ' +
      'się go przerysować po zmianie liczb.',
  },
  {
    kod: 'design.stock-browser',
    tytul: 'Przeglądarka baz zdjęciowych',
    rola: 'zarzadca',
    grupa: 'bazy',
    objasnienie:
      'Wyszukanie w bazach zdjęciowych i wciągnięcie zasobu wraz z licencją. Cztery bazy ' +
      'pracują BEZ KLUCZA i bez zakładania konta — od pierwszej minuty po instalacji. ' +
      'Cztery czytają klucz z sejfu rdzenia i przy jego braku wracają w wykazie ' +
      'dostawców, którzy nie odpowiedzieli — nie w ciszy i nie odmową całej komendy. ' +
      'Licencja zapisuje się RAZEM z zasobem: materiał, o którym nikt później nie powie, ' +
      'czy wolno go było użyć, jest usterką, nie zasobem.',
  },
  {
    kod: 'design.mockup-workshop',
    tytul: 'Warsztat makiety',
    rola: 'pomocnicze',
    grupa: 'makieta',
    objasnienie:
      'Ramki jako ekrany, układ automatyczny, więzy responsywne, siatki, komponenty ' +
      'z wariantami i ich instancje, przejścia prototypu oraz dwie drogi makiety: ' +
      'z opisu i ze zrzutu ekranu. Ramka jest EKRANEM, kompozycja PŁÓTNEM — to rozmiar ' +
      'ramki sprawdza się na innym urządzeniu, a więzy przeliczają wtedy warstwy. ' +
      'Odczyt prototypu mówi, do których ramek NIE prowadzi żadne przejście.',
  },
  {
    kod: 'design.publication-workshop',
    tytul: 'Warsztat publikacji',
    rola: 'pomocnicze',
    grupa: 'publikacja',
    objasnienie:
      'Szablony materiału i publikacje wielostronicowe: strony szablonu, podstawienia ' +
      'treści, komplety kampanii i wydanie. Książka, broszura i katalog powstają jako ' +
      'szablon o wielu stronach i wychodzą JEDNYM plikiem PDF — przez to samo wydanie do ' +
      'druku i przez tę samą kontrolę przeddrukową. Nazwa podstawienia, której szablon ' +
      'nie ma, wraca w wykazie niedopasowanych: cicha zgoda dałaby komplet kampanii, ' +
      'w którym połowa podstawień nie weszła.',
  },
];

/** Pięć okien warsztatowych modułu wraz z ich wspólnym źródłem. */
export interface WarsztatyDesignu {
  okna: readonly OknoWarsztatuDesignu[];
  /** Zleca odczyt materiału we wszystkich warsztatach. */
  wczytaj(): Promise<void>;
}

export function utworzWarsztatyDesignu(stan: StanDesignu, kanal: Kanal): WarsztatyDesignu {
  // Jedno źródło na pięć okien: wszystkie idą tą samą drogą do rdzenia i czytają
  // ten sam magazyn materiału. Pięć źródeł byłoby pięcioma połączeniami do tego
  // samego kanału.
  //
  // Kanał wchodzi wprost, a nie ze stanu modułu: stan oddaje źródła obszaru
  // `design.*` i zaplecza, a warsztaty idą DOWOLNĄ komendą z katalogu czynności.
  // Dołożenie kanału do stanu tylko dla nich otwierałoby wszystkim oknom drogę
  // obok źródeł, które stan dla nich trzyma.
  const zrodlo = utworzZrodloWarsztatowDesignu(kanal);
  const okna = NASTAWY_WARSZTATOW_DESIGNU.map((nastawy) =>
    utworzOknoWarsztatuDesignu(stan, zrodlo, nastawy),
  );
  return {
    okna,
    async wczytaj() {
      // Odczyt idzie równolegle: pięć warsztatów czyta ten sam wykaz, a szeregowe
      // odczyty kazałyby Operatorowi czekać pięć razy na to samo.
      await Promise.all(okna.map((okno) => okno.wczytaj()));
    },
  };
}
