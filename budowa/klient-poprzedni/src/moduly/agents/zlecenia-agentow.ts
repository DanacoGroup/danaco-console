import type {
  AgentConnectorKind,
  AgentPermissionGroup,
  AgentVisibility,
  IdentityLayer,
  MemoryLevel,
  ProviderTransport,
} from '../../../../shared/contract';

/**
 * Zlecenia okien modułu Agents — kształt tego, co okno wie, zanim zbuduje
 * żądanie kontraktu.
 *
 * Okno pracuje polami formularza: napisami z kontrolek i wartościami list
 * wyboru. Kontrakt pracuje polami opcjonalnymi, których pusta wartość ma
 * znaczyć „nie podano”, a nie „podano pustkę”. Te struktury są
 * granicą między jednym a drugim: pole puste zostaje tutaj, a do rdzenia idzie
 * żądanie bez niego. Przekład wykonuje `zrodlo-agentow.ts`.
 */

/**
 * Tożsamość zakładanego eksperta (zakładka „Tożsamość” Agent Buildera).
 *
 * `nazwa` jest nazwą bytu w bibliotece, `imie` — imieniem własnym, po którym
 * Operator poznaje eksperta w wykazie modeli obszaru polecenia. To dwie różne
 * rzeczy, więc dwa pola. `instrukcje` odpowiada `Agent.systemPrompt` i stoi
 * obok warstw — patrz nota w `warstwy-promptu.ts`.
 */
export interface TozsamoscEksperta {
  nazwa: string;
  opis: string;
  instrukcje: string;
  kanal: string;
  model: string;
  imie: string;
  favikon: string;
  /** Odstępstwo od globalnego promptu systemowego — patrz nota niżej. */
  zastepuje: boolean;
  /** Zasięg widoczności eksperta w bibliotece. */
  widocznosc: AgentVisibility;
  /** Poziomy pamięci — patrz nota o zbiorze pustym niżej. */
  poziomyPamieci: readonly MemoryLevel[];
}

/** Zmiana tożsamości — pole pominięte zostaje takie, jakie było. */
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
 * Poziomy pamięci jadą zawsze, także jako zbiór pusty.
 *
 * Kontrakt rozstrzyga to wprost: pominięcie pola znaczy „bez zmiany”, a LISTA
 * PUSTA znaczy „pamięć wyłączona w całości” — i jest jedynym zapisem
 * wyłączenia, bo `MemoryLevel` piątej wartości nie ma. Gdyby okno pomijało pole
 * przy wszystkich poziomach odznaczonych, Operator odznaczyłby cztery pola
 * i nie wyłączyłby niczego.
 */

/**
 * Odstępstwo jest wartością logiczną, a nie `IdentityMode`, bo dwie wartości
 * `Agent.mode` nie są równorzędne: prompt systemowy ustawiany globalnie w oknie
 * konfiguracji obowiązuje domyślnie, a moduł Agents daje albo instrukcję
 * dopisywaną do niego (stan domyślny, `false`), albo jawne odstępstwo (`true`),
 * po którym instrukcja eksperta staje się promptem systemowym. Przekład na
 * `mode` robi `zrodlo-agentow.ts` — na granicy formularza i kontraktu.
 */

/**
 * Zlecenie zapisu jednej warstwy promptu (`agent.layer.set`).
 *
 * Zlecenie nie niesie trybu podania: `AgentLayerSetRequest` pola `mode` nie ma,
 * bo o dopisaniu do globalnego promptu systemowego albo jego zastąpieniu
 * rozstrzyga `Agent.mode`, oznaczany raz przy tożsamości eksperta.
 */
export interface ZlecenieWarstwy {
  idEksperta: string;
  warstwa: IdentityLayer;
  tresc: string;
  czynna: boolean;
}

/** Zlecenie przypisania wtyczki (`agent.plugin.add`). */
export interface ZlecenieWtyczki {
  idEksperta: string;
  nazwa: string;
  zrodlo: string;
  wersja: string;
}

/** Zlecenie okna Model Configuration. */
export interface ZlecenieModelu {
  idEksperta: string;
  kanal: string;
  model: string;
  transport: ProviderTransport | '';
  parametry: string;
}

/** Zlecenie okna Connectors Manager. */
export interface ZlecenieKonektora {
  idEksperta: string;
  nazwa: string;
  rodzaj: AgentConnectorKind;
  idPunktu: string;
  konfiguracja: string;
}

/** Zlecenie okna Permissions Center. */
export interface ZlecenieUprawnienia {
  idEksperta: string;
  grupa: AgentPermissionGroup;
  przyznane: boolean;
  zakres: string;
}
