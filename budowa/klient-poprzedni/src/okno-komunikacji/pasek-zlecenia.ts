import './pasek-zlecenia.css';

import { EventType } from '../../../shared/contract';
import { utworzZrodloRozszerzen } from '../moduly/agents/zrodlo-rozszerzen';
import type { Kanal } from '../protokol/kanal';
import { utworzRejestrAgentow } from '../sterowanie/rejestr-agentow';
import type { RejestrKanalow } from '../sterowanie/rejestr-kanalow';
import { utworzDyktowanie, type Dyktowanie } from './dyktowanie/dyktowanie';
import { utworzSterKatalogow } from './ster-katalogow';
import { utworzSterMikrofonu } from './ster-mikrofonu';
import { utworzSterModelu } from './ster-modelu';
import { utworzSterNakladu } from './ster-nakladu';
import type { SterPaska } from './ster-nastawy';
import { utworzSterRozszerzen } from './ster-rozszerzen';
import { utworzSterUprawnien } from './ster-uprawnien';
import type { ZrodloZlecenia } from './zrodlo-zlecenia';

/**
 * Pasek zlecenia — powierzchnia sterowania kopertą zlecenia, pod czatem.
 *
 * Na tym pasku ustawia się całe otoczenie zadania, zanim treść pójdzie do
 * modelu: gdzie się wykona, na czym, jakim modelem, z jakim wysiłkiem i w jakim
 * trybie zatwierdzania. To jedna decyzja rozłożona na kilka nastaw, więc stoi
 * w jednym rzędzie — dwa rzędy byłyby dwiema kopertami.
 *
 * W rzędzie stoją: urządzenie (`przelacznik-srodowiska.ts`, wchodzi z zewnątrz
 * jako gotowy element), katalog roboczy, model, wysiłek, tryb zatwierdzania
 * i rozszerzenia.
 *
 * Mikrofon jest zbudowany, ale nie staje: `utworzSterMikrofonu` oddaje `null`,
 * bo dostępność dyktowania wymaga obu warunków naraz, a drugi jest
 * niespełniony — kontrakt nie ma komendy przenoszącej nagranie z pamięci karty
 * na maszynę silnika (`dyktowanie/dostarczenie-nagrania.ts`, `dostepne()`).
 * Pasek jest wtedy krótszy: wyszarzona ikona byłaby bramą, a ikona odmawiająca
 * po naciśnięciu — atrapą.
 *
 * Nie mają portu w ogóle: projekt (rdzeń nie zna komendy dającej wykaz
 * projektów), izolacja i drzewo równoległe (`isolation.*` to poziomy zasięgu),
 * wersjonowanie i zatwierdzanie zmian (`agent.version.*` należą do koperty
 * Agent Buildera, nie okna czatu), załączniki (`message.send.attachments`
 * przyjmuje odwołania, a komendy wnoszącej plik do okna czatu nie ma) oraz
 * sposób podania (brak komendy dokładającej pozycję do kolejki). Pustych
 * uchwytów na nie tu nie ma.
 *
 * Dwa ostatnie stery wchodzą inaczej niż pozostałe. Stery koperty czytają
 * migawkę okna i nic więcej nie potrzebują. Rozszerzenia mają własny katalog
 * w rdzeniu i własne zdarzenie, a mikrofon — własny sprzęt; oba trzymają
 * subskrypcje, których pasek pozbywa się przy zejściu okna (`rozlacz`).
 *
 * Mikrofon staje dopiero po zapytaniu o silnik: `Dyktowanie.dostepnosc()` pada,
 * zanim ikona powstanie, więc ster wchodzi do rzędu z opóźnieniem jednej
 * odpowiedzi rdzenia. Miejsce ma przygotowane z góry (gniazdo), żeby kolejność
 * rzędu nie zależała od czasu odpowiedzi.
 *
 * Pasek nie buduje menu i nie zna komend: menu buduje
 * `komponenty/menu-drzewo.ts`, wysyłkę prowadzi `zrodlo-zlecenia.ts`. Ten plik
 * ustawia stery w rzędzie i podaje każdemu wycinek migawki, który do niego
 * należy.
 *
 * Wszystkie stery czytają jedną migawkę okna i piszą jedną drogą — tymi samymi
 * bytami `sterowanie/`, którymi pisze kolumna sterowania.
 */

/** Zależności paska — wąskie i wstrzykiwane. */
export interface ZaleznosciPaska {
  /** Droga do rdzenia — wyłącznie po to, by założyć rejestr ekspertów. */
  kanal: Kanal;
  /** Odczyt i zapis nastaw okna; jeden na pasek. */
  zrodlo: ZrodloZlecenia;
  /** Wspólny katalog kanałów modelu — jeden na klienta. */
  rejestrKanalow: RejestrKanalow;
  /** Gotowy przełącznik urządzenia; pominięty = paska bez tego steru. */
  sterowanieSrodowiska?: HTMLElement;
  /** Otwiera kolumnę sterowania okna — stopka steru katalogów. */
  otworzSterowanie(): void;
  /**
   * Dokłada rozpoznany tekst do pola wypowiedzi; warunek postawienia mikrofonu.
   *
   * Pominięta znaczy „pasek bez mikrofonu": dyktowanie, którego wynik nie ma
   * dokąd trafić, byłoby przyciskiem meldującym pracę bez skutku.
   */
  wstawTekst?(tekst: string): void;
  /** Otwiera katalog rozszerzeń; pominięta = menu „+" bez stopki. */
  otworzKatalogRozszerzen?(): void;
}

export interface PasekZlecenia {
  /** Element stawiany w rzędzie akcji pod polem wypowiedzi. */
  element: HTMLElement;
  /** Przerysowuje wszystkie stery z bieżącej migawki. */
  odswiez(): void;
  /**
   * Zdejmuje nasłuchy sterów i zwalnia mikrofon; obowiązkowe przy zejściu okna.
   *
   * Nasłuch sprzętu dźwiękowego (`devicechange`) i subskrypcja
   * `extension.changed` żyją poza oknem: bez tego wywołania przeżyłyby je
   * i przerysowywałyby stery zdjęte już z dokumentu, a niezwolniony strumień
   * zostawiłby zapaloną lampkę mikrofonu.
   */
  rozlacz(): void;
}

export function utworzPasekZlecenia(zaleznosci: ZaleznosciPaska): PasekZlecenia {
  const { kanal, zrodlo, rejestrKanalow } = zaleznosci;

  const element = document.createElement('div');
  element.className = 'dc-pasek-zlecenia';
  // Rola `group` niesie tu znaczenie: czytnik ekranu ma powiedzieć, czego
  // dotyczy ten rząd uchwytów, zanim zacznie czytać same wartości.
  element.setAttribute('role', 'group');
  element.setAttribute('aria-label', 'Koperta zlecenia');

  const model = utworzSterModelu(
    {
      migawka: () => {
        const okno = zrodlo.migawka().okno;
        return okno.agentId === undefined
          ? { modelChannelId: okno.modelChannelId }
          : { modelChannelId: okno.modelChannelId, agentId: okno.agentId };
      },
      zastosuj: (nazwa, zmiana) => zrodlo.zastosuj(nazwa, zmiana),
    },
    rejestrKanalow,
    utworzRejestrAgentow(kanal),
  );

  const naklad = utworzSterNakladu({
    migawka: () => ({ naklad: zrodlo.migawka().ustawienia.nakladRozumowania }),
    zapisz: (nazwa, klucz, wartosc) => zrodlo.zapisz(nazwa, klucz, wartosc),
  });

  const uprawnienia = utworzSterUprawnien({
    migawka: () => ({ tryb: zrodlo.migawka().okno.permissionMode }),
    zastosuj: (nazwa, zmiana) => zrodlo.zastosuj(nazwa, zmiana),
  });

  const katalogi = utworzSterKatalogow({
    migawka: () => ({ katalogi: zrodlo.migawka().okno.workingDirs }),
    zastosuj: (nazwa, zmiana) => zrodlo.zastosuj(nazwa, zmiana),
    otworzSterowanie: () => zaleznosci.otworzSterowanie(),
  });

  const rozszerzenia = utworzSterRozszerzen({
    zrodlo: utworzZrodloRozszerzen(kanal),
    // Katalog rozszerzeń zmienia się także poza paskiem — z modułu Agents,
    // z Connectors Managera i z drugiego urządzenia Operatora. Rdzeń rozgłasza
    // `extension.changed` do wszystkich połączeń konta, więc pasek nie musi
    // niczego odpytywać cyklicznie ani zgadywać.
    naZmiane: (sluchacz) => kanal.naZdarzenie(EventType.ExtensionChanged, () => sluchacz()),
    ...(zaleznosci.otworzKatalogRozszerzen === undefined
      ? {}
      : { otworzKatalog: zaleznosci.otworzKatalogRozszerzen }),
  });

  // Kolejność niesie treść, nie porządek alfabetyczny: od tego, gdzie zlecenie
  // się wykona i na czym, przez to, czym je wykona, po to, jak ma pytać o zgodę,
  // a na końcu — czym podaje się treść.
  if (zaleznosci.sterowanieSrodowiska !== undefined) {
    element.append(zaleznosci.sterowanieSrodowiska);
  }
  element.append(
    katalogi.element,
    model.element,
    naklad.element,
    uprawnienia.element,
    rozszerzenia.element,
  );

  // Gniazdo mikrofonu stoi od razu, sam ster wchodzi później. Bez tego pustego
  // miejsca ster wskoczyłby na koniec rzędu albo przed niego — zależnie od tego,
  // kiedy rdzeń odpowie na `speech.availability.get`.
  const gniazdoMikrofonu = document.createElement('span');
  gniazdoMikrofonu.className = 'dc-pasek-zlecenia__gniazdo';
  element.append(gniazdoMikrofonu);

  /** Silnik dyktowania zakładany tylko wtedy, gdy wynik ma dokąd trafić. */
  let dyktowanie: Dyktowanie | null = null;
  let mikrofon: (SterPaska & { rozlacz(): void }) | null = null;

  const wstawTekst = zaleznosci.wstawTekst;
  if (wstawTekst !== undefined) {
    dyktowanie = utworzDyktowanie(kanal);
    void utworzSterMikrofonu({ dyktowanie, wstawTekst }).then((ster) => {
      // `null` znaczy „silnika mowy nie ma" — pasek zostaje krótszy zamiast
      // pokazywać wyszarzoną ikonę.
      if (ster === null) return;
      mikrofon = ster;
      gniazdoMikrofonu.append(ster.element);
    });
  }

  /**
   * Przerysowanie idzie z jednej subskrypcji migawki, nie z osobnej na każdy
   * ster: stery koperty czytają ten sam stan, a osobne nasłuchy byłyby wieloma
   * drogami do tej samej roboty i wieloma miejscami do odsubskrybowania.
   * Rozszerzenia i mikrofon migawki okna nie czytają — ich stan leży w katalogu
   * rdzenia i w sprzęcie — więc budzą się własnymi drogami.
   */
  function odswiez(): void {
    for (const ster of [katalogi, model, naklad, uprawnienia]) ster.odswiez();
  }

  // Subskrypcja migawki nie wraca uchwytem do zdjęcia — trzyma ją
  // `zrodlo-zlecenia.ts` i zdejmuje jego własne `rozlacz()`, wołane przez
  // powłokę tuż obok tego. Pasek zdejmuje wyłącznie to, co sam założył.
  zrodlo.naZmiane(odswiez);

  return {
    element,
    odswiez,

    rozlacz() {
      rozszerzenia.rozlacz();
      mikrofon?.rozlacz();
      // Fasada dyktowania zdejmowana zawsze, także gdy ster nie powstał:
      // nasłuch `devicechange` zakłada się przy jej utworzeniu, a nie przy
      // narysowaniu ikony, więc brak silnika mowy nie zwalnia z posprzątania.
      dyktowanie?.rozlacz();
    },
  };
}
