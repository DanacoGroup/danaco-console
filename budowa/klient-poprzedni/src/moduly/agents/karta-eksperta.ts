import { AgentVisibility, type Agent } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { KOLEJNOSC_POZIOMOW } from './poziomy-pamieci';
import { czyZastepuje } from './warstwy-promptu';

/**
 * Karta eksperta w Bibliotece — jedna zapisana definicja jako kafel siatki.
 *
 * Karta zastępuje wiersz z samą nazwą, bo wybór eksperta jest decyzją
 * podejmowaną na podstawie tego, czym ekspert JEST: jakim modelem mówi, jak
 * szeroko jest widziany, ile ma narzędzi i czy odstępuje od globalnego promptu
 * systemowego. Wiersz z nazwą kazał Operatorowi wejść w edytor, żeby to
 * sprawdzić, i wyjść, gdy trafił nie na tego.
 *
 * Wszystkie wartości karty pochodzą z pól bytu `Agent` oddanych przez rdzeń
 * przy odczycie biblioteki. Karta nie woła ani jednej komendy i niczego nie
 * dolicza — liczba narzędzi ma własny byt (`licznik-narzedzi.ts`) i stoi
 * w oknach eksperta, nie tutaj, żeby nie było dwóch rachunków jednej rzeczy.
 *
 * Stan nie jest tu samym kolorem: każda plakietka niesie napis, a odstępstwo od
 * promptu globalnego dostaje osobny znacznik słowny — zgodnie z regułą, że
 * żeton barwy nie zwalnia komponentu z etykiety.
 */

/** Czynności karty — każda dotyczy tego jednego eksperta. */
export interface DzialaniaKarty {
  naWybor(): void;
  naDuplikowanie(): void;
  naUsuniecie(): void;
}

/** Napis plakietki zasięgu widoczności w języku Operatora. */
function napisWidocznosci(zasieg: AgentVisibility): string {
  return zasieg === AgentVisibility.Project ? 'projektowy' : 'globalny';
}

/** Napis plakietki pamięci; zbiór pusty jest wyłączeniem, nie brakiem danych. */
function napisPamieci(ekspert: Agent): string {
  const ile = ekspert.memoryLevels.length;
  if (ile === 0) return 'pamięć wyłączona';
  if (ile === KOLEJNOSC_POZIOMOW.length) return 'pamięć pełna';
  return `pamięć ${ile}/${KOLEJNOSC_POZIOMOW.length}`;
}

/** Wiersz „model · kanał” — skrót konfiguracji z okna Model Configuration. */
function napisModelu(ekspert: Agent): string {
  const model = ekspert.model ?? '';
  const kanal = ekspert.channelId ?? '';
  if (model === '' && kanal === '') return 'model bazowy nieustalony';
  return [model === '' ? 'model kanału' : model, kanal].filter((czesc) => czesc !== '').join(' · ');
}

/** Plakietka z napisem — barwę niesie znacznik `data-`, arkusz modułu ją maluje. */
function plakietka(napis: string, rodzaj: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dn-plakietka da-karta__plakietka';
  element.dataset['rodzaj'] = rodzaj;
  element.textContent = napis;
  return element;
}

export function utworzKarteEksperta(
  ekspert: Agent,
  czynny: boolean,
  dzialania: DzialaniaKarty,
  przypisan?: number,
): HTMLElement {
  // Favikon pełni rolę awatara kwadratowego zarezerwowanego dla agentów.
  // Ekspert bez znaku dostaje pierwszą literę nazwy, a nie pustkę — kafel bez
  // niczego w lewym górnym rogu wygląda na niedoczytany.
  const awatar = document.createElement('span');
  awatar.className = 'dn-awatar dn-awatar--kwadrat da-karta__awatar';
  const znak = (ekspert.favicon ?? '').trim();
  awatar.textContent = znak === '' ? ekspert.name.slice(0, 1).toLocaleUpperCase('pl') : znak;
  awatar.dataset['zastepczy'] = String(znak === '');

  const nazwa = document.createElement('span');
  nazwa.className = 'da-karta__nazwa';
  const imie = (ekspert.displayName ?? '').trim();
  nazwa.textContent = imie === '' ? ekspert.name : `${ekspert.name} — „${imie}”`;

  const wybor = document.createElement('button');
  wybor.type = 'button';
  wybor.className = 'da-karta__wybor';
  wybor.setAttribute('aria-pressed', String(czynny));
  wybor.append(awatar, nazwa);
  wybor.addEventListener('click', dzialania.naWybor);

  const opis = document.createElement('p');
  opis.className = 'da-karta__opis';
  const przeznaczenie = (ekspert.description ?? '').trim();
  opis.textContent = przeznaczenie === '' ? 'bez opisu przeznaczenia' : przeznaczenie;
  opis.dataset['pusty'] = String(przeznaczenie === '');

  const model = document.createElement('p');
  model.className = 'da-karta__model';
  model.textContent = napisModelu(ekspert);

  const plakietki = document.createElement('div');
  plakietki.className = 'da-karta__plakietki';
  plakietki.append(
    plakietka(`wersja ${ekspert.version ?? 1}`, 'wersja'),
    plakietka(napisWidocznosci(ekspert.visibility), 'widocznosc'),
    plakietka(ekspert.enabled ? 'czynny' : 'wyłączony', ekspert.enabled ? 'czynny' : 'wylaczony'),
    plakietka(napisPamieci(ekspert), 'pamiec'),
  );
  // Liczba przypisań pochodzi z `agent.assignment.list` — odczytu OD STRONY
  // EKSPERTA. `workspace.agent.assign` zapisuje przynależność, ale patrzy na nią
  // od strony projektu, więc bez tamtej komendy karta nie miałaby skąd wziąć
  // liczby. Wartość nieznana nie daje plakietki: zero i „nie pytaliśmy" to dwie
  // różne odpowiedzi i karta nie ma prawa ich mylić.
  if (przypisan !== undefined) {
    const znacznik = plakietka(
      przypisan === 0 ? 'bez przypisań' : `${przypisan} przypisań`,
      'przypisania',
    );
    znacznik.title =
      'Projekty modułu Workspace i role środowiska MultitaskingAI, w których ten ekspert ' +
      'jest wykorzystywany jako wykonawca.';
    plakietki.append(znacznik);
  }
  if (czyZastepuje(ekspert)) {
    const znacznik = plakietka('zastępuje prompt globalny', 'odstepstwo');
    znacznik.title =
      'Ekspert ma oznaczone odstępstwo od ustawień domyślnych: jego instrukcja jest ' +
      'promptem systemowym, a globalny prompt z okna konfiguracji i ustawień nie ' +
      'obowiązuje dla jego okien.';
    plakietki.append(znacznik);
  }

  const akcje = document.createElement('div');
  akcje.className = 'da-karta__akcje';
  akcje.append(
    dzialanie('Duplikuj', dzialania.naDuplikowanie),
    dzialanie('Usuń', dzialania.naUsuniecie),
  );

  const element = document.createElement('li');
  element.className = 'dn-karta da-karta';
  element.dataset['ekspert'] = ekspert.id;
  element.dataset['czynny'] = String(czynny);
  element.append(wybor, opis, model, plakietki, akcje);
  return element;
}

function dzialanie(tresc: string, naKlik: () => void): HTMLButtonElement {
  const element = przycisk(tresc, 'dn-btn dn-btn--sm dn-btn--zarys');
  element.addEventListener('click', naKlik);
  return element;
}
