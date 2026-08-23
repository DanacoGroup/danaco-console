import './sterowanie.css';

import type { Window } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { utworzSterowanieHostu } from './host-wykonania';
import { utworzSterowanieKatalogow } from './katalogi-robocze';
import { utworzSterowanieModelu } from './model-glowny';
import { utworzModelKartySesji } from './model-karty-sesji';
import { utworzSterowanieModeluZapasowego } from './model-zapasowy';
import { utworzSterowanieModulu } from './modul-okna';
import { utworzRejestrModulow } from './rejestr-modulow';
import { utworzSterowanieNakladu } from './naklad-rozumowania';
import { utworzPasekKomunikatow } from './pasek-komunikatow';
import { utworzRejestrAgentow } from './rejestr-agentow';
import type { RejestrKanalow } from './rejestr-kanalow';
import { utworzSterowanieRoli } from './rola-okna';
import { utworzSterowanieSrodowiska } from './srodowisko-wykonania';
import { utworzStanSterowania } from './stan-sterowania';
import { utworzSterowanieUprawnien } from './tryb-uprawnien';
import { utworzZmianeOkna } from './zmiana-okna';
import { utworzZmianeUstawienia } from './zmiana-ustawienia';

/** Zależności kompletu: kanał kontraktu, okno oraz wspólny wykaz kanałów modelu. */
export interface ZaleznosciPanelu {
  kanal: Kanal;
  okno: Window;
  rejestrKanalow: RejestrKanalow;
}

/** Komplet sterowania jednego okna komunikacji. */
export interface PanelSterowania {
  /** Element montowany przy oknie. */
  element: HTMLElement;
  /** Okno, którego dotyczy komplet. */
  idOkna(): string;
  /** Przyjmuje okno potwierdzone przez rdzeń poza kompletem. */
  przyjmijOkno(okno: Window): void;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

/**
 * Komplet sterowania **per okno**, nie globalny pasek.
 *
 * Każde wywołanie buduje osobny stan, osobne sterowania i osobne subskrypcje,
 * domknięte na identyfikatorze swojego okna. Dwa komplety otwarte obok siebie
 * nie mają wspólnej zmiennej: zmiana w jednym idzie komendą `window.update`
 * z identyfikatorem tego okna, a potwierdzenie — odpowiedź i zdarzenie
 * `window.changed` — trafia wyłącznie do stanu okna o tym identyfikatorze.
 *
 * Wspólny pozostaje jedynie wykaz kanałów modelu: katalog wyboru, nie
 * ustawienie okna.
 */
export function utworzPanelSterowania(zaleznosci: ZaleznosciPanelu): PanelSterowania {
  const { kanal, okno, rejestrKanalow } = zaleznosci;

  const stan = utworzStanSterowania(okno);
  const komunikaty = utworzPasekKomunikatow();
  const zmiana = utworzZmianeOkna(kanal, stan, komunikaty.pokaz);
  const ustawienia = utworzZmianeUstawienia(kanal, stan, komunikaty.pokaz);

  const element = document.createElement('section');
  element.className = 'dc-sterowanie';
  element.dataset.okno = stan.idOkna();
  element.append(
    naglowek(stan.idOkna()),
    siatka([
      utworzSterowanieSrodowiska(stan, zmiana),
      utworzSterowanieHostu(stan, ustawienia),
      utworzSterowanieModulu(stan, zmiana, utworzRejestrModulow(kanal)),
      // Rejestr ekspertów powstaje w miejscu, wzorem rejestru modułów wiersz
      // wyżej: jest katalogiem wyboru czytanym z rdzenia, nie stanem okna,
      // więc nie ma po co przeciągać go przez umowę kompletu.
      // Rozesłanie wyboru na całą kartę sesji idzie osobną komendą
      // (`model.channel.set` z `sessionId`) — to jedyna jej zdolność, której
      // `window.update` nie ma.
      utworzSterowanieModelu(
        stan,
        zmiana,
        rejestrKanalow,
        utworzRejestrAgentow(kanal),
        utworzModelKartySesji(stan, kanal, komunikaty.pokaz),
      ),
      utworzSterowanieModeluZapasowego(stan, ustawienia, rejestrKanalow),
      utworzSterowanieNakladu(stan, ustawienia),
      utworzSterowanieUprawnien(stan, zmiana),
      utworzSterowanieRoli(stan, zmiana),
      utworzSterowanieKatalogow(stan, zmiana),
    ]),
    komunikaty.element,
  );

  ustawienia.wczytaj();

  return {
    element,
    idOkna: stan.idOkna,
    przyjmijOkno: stan.przyjmijOkno,
    rozlacz() {
      zmiana.rozlacz();
      ustawienia.rozlacz();
    },
  };
}

/** Nagłówek kompletu z identyfikatorem okna — komplet należy do jednego okna. */
function naglowek(idOkna: string): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dc-sterowanie__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dc-sterowanie__tytul';
  tytul.textContent = 'Sterowanie okna';

  const identyfikator = document.createElement('span');
  identyfikator.className = 'dc-sterowanie__okno';
  identyfikator.textContent = idOkna;
  identyfikator.title = `Identyfikator okna komunikacji: ${idOkna}`;

  element.append(tytul, identyfikator);
  return element;
}

/** Siatka sterowań kompletu. */
function siatka(sterowania: HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dc-sterowanie__siatka';
  element.append(...sterowania);
  return element;
}
