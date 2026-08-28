import { AgentVisibility, type Agent } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { KOLEJNOSC_POZIOMOW } from './poziomy-pamieci';
import { czyZastepuje } from './warstwy-promptu';

/**
 * Karta eksperta w Bibliotece pokazuje jedną zapisaną definicję jako kafel
 * siatki, złożony wyłącznie z pól bytu oddanych przez rdzeń.
 */
/** Czynności karty, z których każda dotyczy tego jednego eksperta i jego pozycji w bibliotece ekspertów. */
export interface DzialaniaKarty {
  naWybor(): void;
  naDuplikowanie(): void;
  naUsuniecie(): void;
}

/** Napis plakietki zasięgu widoczności eksperta na karcie, wyrażony wprost w języku zrozumiałym dla Operatora. */
function napisWidocznosci(zasieg: AgentVisibility): string {
  return zasieg === AgentVisibility.Project ? 'projektowy' : 'globalny';
}

/** Napis plakietki pamięci eksperta na karcie; zbiór pusty jest wyłączeniem pamięci, nie brakiem danych. */
function napisPamieci(ekspert: Agent): string {
  const ile = ekspert.memoryLevels.length;
  if (ile === 0) return 'pamięć wyłączona';
  if (ile === KOLEJNOSC_POZIOMOW.length) return 'pamięć pełna';
  return `pamięć ${ile}/${KOLEJNOSC_POZIOMOW.length}`;
}

/** Wiersz „model · kanał” widoczny na karcie eksperta jako skrót konfiguracji z okna Model Configuration. */
function napisModelu(ekspert: Agent): string {
  const model = ekspert.model ?? '';
  const kanal = ekspert.channelId ?? '';
  if (model === '' && kanal === '') return 'model bazowy nieustalony';
  return [model === '' ? 'model kanału' : model, kanal].filter((czesc) => czesc !== '').join(' · ');
}

/** Plakietka z napisem na karcie, której barwę niesie znacznik danych, a maluje ją arkusz stylu modułu. */
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
  // Favikon pełni rolę awatara; ekspert bez znaku dostaje pierwszą literę nazwy, nie pustkę.
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
  // Liczba przypisań pochodzi z odczytu od strony eksperta; wartość nieznana nie daje plakietki.
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
