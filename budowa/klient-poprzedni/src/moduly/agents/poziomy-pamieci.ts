import { MemoryLevel } from '../../../../shared/contract';

/**
 * Konfiguracja pamięci eksperta wybiera niezależnie cztery poziomy zasięgu,
 * ponieważ ekspert korzysta z kilku naraz. Wyłączeniem pamięci jest zbiór
 * pusty, czyli cztery pola odznaczone; piątego poziomu kontrakt nie zna.
 */
export interface PoziomyPamieci {
  /** Grupa osadzana w formularzu tożsamości. */
  element: HTMLElement;
  /** Nanosi poziomy eksperta; `null` znaczy formularz zakładania. */
  ustaw(poziomy: readonly MemoryLevel[] | null): void;
  /** Poziomy zaznaczone w chwili odczytu — zbiór pusty znaczy pamięć wyłączoną. */
  wybrane(): readonly MemoryLevel[];
}

/**
 * Kolejność poziomów biegnie od najszerszego do najwęższego, czyli tak, jak
 * schodzi zasięg pamięci: od pamięci wspólnej wszystkim projektom po pamięć
 * jednej karty rozmowy. Grupa kontrolek rysuje się w tym samym porządku.
 */
export const KOLEJNOSC_POZIOMOW: readonly MemoryLevel[] = [
  MemoryLevel.Global,
  MemoryLevel.Project,
  MemoryLevel.Environment,
  MemoryLevel.Session,
];

/**
 * Nazwy poziomów w języku Operatora wraz ze zdaniem o zasięgu każdego z nich.
 * Nazwy stoją osobno od wartości kontraktu, ponieważ Operator czyta zasięg
 * pamięci, a rdzeń przyjmuje wartości wyliczenia `MemoryLevel`.
 */
export const NAZWY_POZIOMOW: Record<MemoryLevel, string> = {
  [MemoryLevel.Global]: 'globalna — wspólna wszystkim projektom i oknom',
  [MemoryLevel.Project]: 'projekt — pamięć projektu, do którego należy okno',
  [MemoryLevel.Environment]: 'środowisko — TalkIn, WorkSpace, CodeStudio, MultitaskingAI',
  [MemoryLevel.Session]: 'sesja — pamięć jednej karty rozmowy',
};

/**
 * Poziomy wyjściowe nowego eksperta obejmują komplet czterech, ponieważ tak
 * stan wyjściowy nazywa kontrakt. Ruszanie z pustki zakładałoby eksperta bez
 * pamięci przy każdym zapisie, w którym Operator grupy nie dotknął.
 */
const POZIOMY_WYJSCIOWE: readonly MemoryLevel[] = KOLEJNOSC_POZIOMOW;

/**
 * Zdanie pod grupą czyta zaznaczenia bieżące, a nie zamiar Operatora, i zmienia
 * się przy każdym kliknięciu. Dzięki temu Operator widzi, że odznaczenie
 * wszystkiego jest wyłączeniem pamięci, zanim naciśnie zapis.
 */
function zdanieOStanie(wybrane: readonly MemoryLevel[]): string {
  if (wybrane.length === 0) {
    return (
      'Żaden poziom nie jest zaznaczony — zapis WYŁĄCZY pamięć tego eksperta w całości. ' +
      'Wyłączenie nie ma własnej pozycji na liście: zbiór pusty jest jego jedynym zapisem.'
    );
  }
  if (wybrane.length === KOLEJNOSC_POZIOMOW.length) {
    return 'Wszystkie cztery poziomy czynne — stan wyjściowy definicji eksperta.';
  }
  return `Poziomów czynnych: ${wybrane.length} z ${KOLEJNOSC_POZIOMOW.length}.`;
}

export function utworzPoziomyPamieci(): PoziomyPamieci {
  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Konfiguracja pamięci';

  const nota = document.createElement('p');
  nota.className = 'dn-pole-opis';
  nota.textContent =
    'Poziomy, z których ekspert korzysta domyślnie. Wartość definicji jest nadpisywana ' +
    'przy przypisaniu do projektu i przy przypisaniu do roli — pierwszeństwo ma zasięg ' +
    'najbardziej szczegółowy.';

  const stan = document.createElement('p');
  stan.className = 'dn-pole-opis da-pamiec__stan';
  stan.setAttribute('role', 'note');

  const lista = document.createElement('ul');
  lista.className = 'da-pamiec';

  const kontrolki = new Map<MemoryLevel, HTMLInputElement>();

  for (const poziom of KOLEJNOSC_POZIOMOW) {
    const kontrolka = document.createElement('input');
    kontrolka.type = 'checkbox';
    kontrolka.className = 'dn-check';
    kontrolka.id = `da-pamiec-${poziom}`;
    kontrolka.addEventListener('change', () => odswiezZdanie());

    const etykieta = document.createElement('label');
    etykieta.className = 'dn-pole-etykieta';
    etykieta.htmlFor = kontrolka.id;
    etykieta.textContent = NAZWY_POZIOMOW[poziom];

    const wiersz = document.createElement('li');
    wiersz.className = 'da-pamiec__wiersz';
    wiersz.dataset['poziom'] = poziom;
    wiersz.append(kontrolka, etykieta);

    kontrolki.set(poziom, kontrolka);
    lista.append(wiersz);
  }

  const element = document.createElement('section');
  element.className = 'da-panel da-pamiec-panel';
  element.append(tytul, nota, lista, stan);

  function wybrane(): readonly MemoryLevel[] {
    return KOLEJNOSC_POZIOMOW.filter((poziom) => kontrolki.get(poziom)?.checked === true);
  }

  function odswiezZdanie(): void {
    const czynne = wybrane();
    stan.textContent = zdanieOStanie(czynne);
    stan.dataset['wylaczona'] = String(czynne.length === 0);
  }

  function ustaw(poziomy: readonly MemoryLevel[] | null): void {
    const zrodlo = poziomy ?? POZIOMY_WYJSCIOWE;
    for (const [poziom, kontrolka] of kontrolki) kontrolka.checked = zrodlo.includes(poziom);
    odswiezZdanie();
  }

  ustaw(null);

  return { element, ustaw, wybrane };
}
