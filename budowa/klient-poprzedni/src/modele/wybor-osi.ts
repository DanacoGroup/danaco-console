import { ConfigAxis } from '../../../shared/contract';
import { BYTY_OSI, NAZWY_OSI, OSIE_OD_NAJWEZSZEJ, osWymagaBytu } from '../konfiguracja/zasiegi';
import { poleWyboru } from './kontrolki-formularza';

/**
 * Wskazanie osi rozstrzygania wraz z bytem stanowi wspólny nagłówek trzech paneli
 * sekcji modeli: ustawień, tożsamości oraz podglądu promptu. Pasek oddaje wybraną
 * oś i identyfikator jej bytu, a każdą zmianę wskazania rozgłasza słuchaczom.
 */
export interface WskazanieOsi {
  /** Oś rozstrzygania. */
  os: ConfigAxis;
  /** Byt osi: identyfikator modelu albo konta; pusty dla osi platformy. */
  bytOsi: string;
}

export interface WyborOsi {
  /** Pasek osadzany nad panelem. */
  element: HTMLElement;
  /** Oś i byt wskazane kontrolkami. */
  wskazanie(): WskazanieOsi;
  /** Nanosi podpowiedzi bytu: identyfikatory modeli oraz kont. */
  ustawPodpowiedzi(modele: readonly string[], konta: readonly PodpowiedzKonta[]): void;
  /** Subskrypcja zmiany wskazania. */
  naZmiane(sluchacz: () => void): void;
}

/**
 * Konto jako podpowiedź bytu osi: identyfikator służący adresowaniu oraz nazwa
 * przeznaczona do wydruku w polu wyboru. Nazwa nie bierze udziału ani
 * w porównaniach, ani w budowaniu żądań kierowanych do rdzenia.
 */
export interface PodpowiedzKonta {
  identyfikator: string;
  nazwa: string;
}

export function utworzWyborOsi(): WyborOsi {
  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  const os = poleWyboru(
    { etykieta: 'Oś' },
    OSIE_OD_NAJWEZSZEJ.map((wartosc) => ({ wartosc, etykieta: NAZWY_OSI[wartosc] })),
  );
  os.kontrolka.value = ConfigAxis.Platform;

  const podpowiedzi = document.createElement('datalist');
  podpowiedzi.id = 'dm-byty-osi';

  const byt = document.createElement('input');
  byt.type = 'text';
  byt.className = 'dn-pole-kontrolka';
  byt.id = 'dm-byt-osi';
  byt.setAttribute('list', podpowiedzi.id);

  const etykietaBytu = document.createElement('label');
  etykietaBytu.className = 'dn-pole-etykieta';
  etykietaBytu.htmlFor = byt.id;
  etykietaBytu.textContent = 'Byt osi';

  const koszykBytu = document.createElement('div');
  koszykBytu.className = 'dn-pole dm-pole';
  koszykBytu.append(etykietaBytu, byt, podpowiedzi);

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'dn-pole-opis dm-os__wyjasnienie';

  const element = document.createElement('div');
  element.className = 'dm-os';
  element.append(os.element, koszykBytu, wyjasnienie);

  /** Podpowiedzi stoją rozdzielone, ponieważ wybór osi rozstrzyga, które są widoczne. */
  let modele: readonly string[] = [];
  let konta: readonly PodpowiedzKonta[] = [];

  function wybranaOs(): ConfigAxis {
    return os.kontrolka.value as ConfigAxis;
  }

  function ubierz(): void {
    const czynna = wybranaOs();
    koszykBytu.hidden = !osWymagaBytu(czynna);
    byt.placeholder = BYTY_OSI[czynna] ?? '';
    wyjasnienie.textContent = WYJASNIENIA[czynna] ?? '';
    podpowiedzi.replaceChildren(
      ...(czynna === ConfigAxis.Account
        ? konta.map((konto) => pozycja(konto.identyfikator, konto.nazwa))
        : modele.map((model) => pozycja(model, model))),
    );
  }

  for (const kontrolka of [os.kontrolka, byt]) {
    kontrolka.addEventListener('change', () => {
      ubierz();
      oglos();
    });
  }

  ubierz();

  return {
    element,

    wskazanie: () => ({
      os: wybranaOs(),
      bytOsi: koszykBytu.hidden ? '' : byt.value.trim(),
    }),

    ustawPodpowiedzi(noweModele, noweKonta) {
      modele = noweModele;
      konta = noweKonta;
      ubierz();
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}

/**
 * Rozstrzyga, czy wskazanie niesie komplet potrzebny do adresowania, to jest oś
 * wraz z jej bytem. Oś platformy bytu nie wymaga, więc dla niej samo wskazanie
 * osi jest już kompletne.
 */
export function wskazanieKompletne(wskazanie: WskazanieOsi): boolean {
  return !osWymagaBytu(wskazanie.os) || wskazanie.bytOsi !== '';
}

function pozycja(wartosc: string, etykieta: string): HTMLOptionElement {
  const element = document.createElement('option');
  element.value = wartosc;
  element.label = etykieta;
  return element;
}

const WYJASNIENIA: Readonly<Record<ConfigAxis, string>> = {
  [ConfigAxis.Platform]:
    'Zapis obowiązuje każdy model i każde konto, o ile oś węższa nie powie inaczej.',
  [ConfigAxis.Model]:
    'Zapis obowiązuje wskazany model niezależnie od konta, którego poświadczeniem się do niego łączymy.',
  [ConfigAxis.Account]:
    'Zapis obowiązuje wskazane konto i wygrywa z zapisem modelu oraz platformy.',
};
