/**
 * Zlecenia okien modułu Agents opisują kształt danych zebranych z pól
 * formularza, zanim okno złoży z nich żądanie kontraktu. Przekład pola pustego
 * na pole pominięte wykonuje `zrodlo-agentow.ts`.
 */
import type {
  AgentConnectorKind,
  AgentPermissionGroup,
  AgentVisibility,
  IdentityLayer,
  MemoryLevel,
  ProviderTransport,
} from '../../../../shared/contract';

/**
 * Tożsamość zakładanego eksperta zbierana w zakładce tożsamości. Pole `nazwa`
 * nazywa byt w bibliotece, a `imie` jest imieniem własnym, po którym Operator
 * poznaje eksperta w wykazie modeli obszaru polecenia.
 */
export interface TozsamoscEksperta {
  nazwa: string;
  opis: string;
  instrukcje: string;
  kanal: string;
  model: string;
  imie: string;
  favikon: string;
  /** Odstępstwo od globalnego promptu systemowego zamiast dopisania do niego. */
  zastepuje: boolean;
  /** Zasięg widoczności eksperta w bibliotece. */
  widocznosc: AgentVisibility;
  /** Poziomy pamięci eksperta; lista pusta wyłącza pamięć w całości. */
  poziomyPamieci: readonly MemoryLevel[];
}

/**
 * Zmiana tożsamości zakładanego eksperta: każde pole jest opcjonalne, a pole
 * pominięte zostaje w rdzeniu takie, jakie było. Zbiór poziomów pamięci jest
 * tu wyjątkiem, ponieważ lista pusta wyłącza pamięć w całości.
 */
export interface ZmianaEksperta {
  nazwa?: string;
  opis?: string;
  instrukcje?: string;
  czynny?: boolean;
  imie?: string;
  favikon?: string;
  zastepuje?: boolean;
  widocznosc?: AgentVisibility;
  poziomyPamieci?: readonly MemoryLevel[];
}

/**
 * Zlecenie zapisu jednej warstwy promptu komendą `agent.layer.set`. Trybu
 * podania nie niesie, ponieważ o dopisaniu do globalnego promptu systemowego
 * albo o jego zastąpieniu rozstrzyga pole `Agent.mode` tożsamości eksperta.
 */
export interface ZlecenieWarstwy {
  idEksperta: string;
  warstwa: IdentityLayer;
  tresc: string;
  czynna: boolean;
}

/**
 * Zlecenie przypisania wtyczki do eksperta komendą `agent.plugin.add`, złożone
 * z identyfikatora eksperta, nazwy wtyczki, jej źródła oraz żądanej wersji.
 */
export interface ZlecenieWtyczki {
  idEksperta: string;
  nazwa: string;
  zrodlo: string;
  wersja: string;
}

/**
 * Zlecenie okna Model Configuration: wybór kanału i modelu dla eksperta wraz
 * z transportem dostawcy oraz parametrami wywołania podanymi jako napis,
 * który rdzeń rozbiera dopiero po swojej stronie.
 */
export interface ZlecenieModelu {
  idEksperta: string;
  kanal: string;
  model: string;
  transport: ProviderTransport | '';
  parametry: string;
}

/**
 * Zlecenie okna Connectors Manager: jedno połączenie eksperta z zasobem
 * zewnętrznym, opisane nazwą własną, rodzajem konektora, identyfikatorem
 * punktu końcowego oraz konfiguracją podaną jako napis.
 */
export interface ZlecenieKonektora {
  idEksperta: string;
  nazwa: string;
  rodzaj: AgentConnectorKind;
  idPunktu: string;
  konfiguracja: string;
}

/**
 * Zlecenie okna Permissions Center: jedna grupa uprawnień eksperta wraz
 * z rozstrzygnięciem, czy jest przyznana, oraz z zakresem zawężającym jej
 * działanie do wskazanych zasobów.
 */
export interface ZlecenieUprawnienia {
  idEksperta: string;
  grupa: AgentPermissionGroup;
  przyznane: boolean;
  zakres: string;
}
