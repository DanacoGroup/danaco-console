import type { RoundtableParticipant, RoundtableStatement } from '../../../../shared/contract';
import {
  utworzKontekstCzytelnosci,
  znacznikMowcy,
  type KontekstCzytelnosci,
} from './czytelnosc-glosow';
import {
  ostatniaWypowiedz,
  wyciszony,
  zdanieBezSlow,
  zlozGlosBiezacy,
  type GlosBiezacy,
} from './glos-biezacy';
import { zdaniaMetryk, zmierzOdpowiedz } from './metryki-odpowiedzi';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';
import type { StanTresci } from './stany-okna';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';

/** Skład debaty w Model Panels: nastawy widoku wnoszone przez okno przy każdym rysowaniu, bo rysowanie jest czystą funkcją danych. */
export interface NastawySkladu {
  /** Tryb ślepy — tożsamość ukryta do zdjęcia nastawy, wyłącznie widokiem tego okna, nie stanem rdzenia. */
  trybSlepy: boolean;
  /** Czy pokazywać metryki odpowiedzi — długość i czas od otwarcia tury. */
  metryki: boolean;
}

/**
 * Rysuje skład debaty jako osobny panel na każdego uczestnika, albo stan
 * pustki, gdy debata jeszcze nie ma nikogo — to jest stan poprawny (świeżo
 * otwarte okno nie zna składu, dopóki nie przyjdzie przyrost zdarzenia),
 * nie usterka.
 */
export function rysujSkladModeli(
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  tresc: StanTresci,
  nastawy: NastawySkladu,
): void {
  const uczestnicy = stan.uczestnicy();
  if (uczestnicy.length === 0) {
    tresc.pusto(
      'Debata nie ma jeszcze uczestników — dodaj pierwszego z rejestru kanałów powyżej. ' +
        '(Świeżo otwarte okno nie zna składu debaty sprzed otwarcia: odczyt roundtable.model.list jest w kontrakcie, ale okno go jeszcze nie wywołuje.)',
    );
    return;
  }
  const wypowiedzi = stan.wypowiedzi();
  const ostatniMowca = mowcaOstatniegoGlosu(wypowiedzi);
  const kontekst = utworzKontekstCzytelnosci();
  const sklad = document.createElement('div');
  sklad.className = 'dr-sklad';
  uczestnicy.forEach((uczestnik, indeks) => {
    sklad.append(
      panelUczestnika(uczestnik, indeks, stan, {
        kontekst,
        strumien,
        wypowiedzi,
        ostatniMowca,
        nastawy,
      }),
    );
  });
  const miejsce = tresc.tresc();
  if (nastawy.trybSlepy) miejsce.append(notaTrybuSlepego());
  miejsce.append(sklad);
}

/**
 * Zdanie o zasięgu trybu ślepego.
 *
 * Bez niego ukryta nazwa czytałaby się jako anonimizacja debaty, a jest
 * anonimizacją jednego widoku: wypowiedzi w Debate Panelu i w zestawieniu
 * udziału nadal noszą tożsamość mówcy.
 */
function notaTrybuSlepego(): HTMLElement {
  const nota = document.createElement('p');
  nota.className = 'dr-sklad__nota';
  nota.dataset['brakSygnalu'] = 'tak';
  nota.textContent =
    'Tryb ślepy ukrywa tożsamości wyłącznie w tym oknie i wyłącznie przed Operatorem. Kontrakt nie ma nastawy trybu ślepego, ' +
    'więc anonimizacja nie schodzi do rdzenia, nie obejmuje Debate Panelu ani zestawienia udziału i nie zmienia tego, ' +
    'co widzą modele uczestniczące w debacie.';
  return nota;
}

/** Zależności rysowania jednego panelu — zebrane, żeby sygnatura funkcji rysującej panel nie puchła kolejnymi parametrami. */
interface OtoczenieSkladu {
  kontekst: KontekstCzytelnosci;
  strumien: StrumienWypowiedzi;
  wypowiedzi: readonly RoundtableStatement[];
  ostatniMowca: string;
  nastawy: NastawySkladu;
}

/** Mówca wypowiedzi stojącej w wykazie tury jako ostatnia; pusty napis, gdy tura nie ma jeszcze żadnej wypowiedzi. */
function mowcaOstatniegoGlosu(wypowiedzi: readonly RoundtableStatement[]): string {
  const ostatnia = wypowiedzi[wypowiedzi.length - 1];
  return ostatnia === undefined ? '' : ostatnia.participantId;
}

/**
 * Panel pojedynczego uczestnika — jego tożsamość, kanał, stan i odpowiedź, znaczony tym samym
 * znacznikiem co w Debate Panelu.
 */
function panelUczestnika(
  uczestnik: RoundtableParticipant,
  indeks: number,
  stan: StanDebaty,
  otoczenie: OtoczenieSkladu,
): HTMLElement {
  const opisKanalu = stan.opisKanalu(uczestnik.channelId);
  const powtorzony = stan.wystapieniaKanalu(uczestnik.channelId) > 1;
  const cisza = wyciszony(uczestnik);
  const biezacy = zlozGlosBiezacy(
    otoczenie.strumien.glos(uczestnik.id),
    ostatniaWypowiedz(otoczenie.wypowiedzi, uczestnik.id),
    cisza,
  );

  const kolejnosc = document.createElement('span');
  kolejnosc.className = 'dr-uczestnik__kolejnosc';
  kolejnosc.textContent = znacznikMowcy(uczestnik.id, stan, otoczenie.kontekst);
  kolejnosc.title = 'Ten sam znacznik stoi nad wypowiedziami tego uczestnika w Debate Panelu.';

  // Tożsamość (nazwa i kanał) siedzi w klasie wspólnej z formularzem dodania uczestnika.
  const slepy = otoczenie.nastawy.trybSlepy;
  const tozsamosc = document.createElement('div');
  tozsamosc.className = 'dr-tozsamosc';
  const nazwa = document.createElement('strong');
  nazwa.textContent = slepy
    ? `Uczestnik ${indeks + 1} — tożsamość ukryta w tym oknie`
    : nazwaUczestnika(uczestnik, opisKanalu);
  tozsamosc.append(nazwa);
  if (!slepy) {
    const kanal = document.createElement('span');
    kanal.textContent = `Kanał: ${opisKanalu === '' ? uczestnik.channelId : opisKanalu}`;
    tozsamosc.append(kanal, promptSystemowy(uczestnik));
  }

  const pola = document.createElement('div');
  pola.className = 'dr-uczestnik__pola';
  pola.append(tozsamosc, listaStanow(uczestnik, indeks, powtorzony, cisza));
  if (otoczenie.nastawy.metryki) {
    pola.append(
      metrykiPanelu(
        biezacy.tekst,
        ostatniaWypowiedz(otoczenie.wypowiedzi, uczestnik.id),
        stan,
      ),
    );
  }

  const akcje = document.createElement('div');
  akcje.className = 'dr-uczestnik__akcje';

  const panel = document.createElement('div');
  panel.className = 'dr-uczestnik';
  // Tożsamość na panelu, nie sam wygląd: po niej odnajduje się uczestnika przy powtórzonym kanale.
  panel.dataset['uczestnik'] = uczestnik.id;
  panel.dataset['kanalPowtorzony'] = powtorzony ? 'tak' : 'nie';
  if (cisza) panel.dataset['wyciszony'] = 'tak';
  if (uczestnik.id === otoczenie.ostatniMowca) panel.dataset['ostatniGlos'] = 'tak';
  panel.append(kolejnosc, pola, akcje, odpowiedzUczestnika(biezacy, cisza));
  return panel;
}

/**
 * Stany uczestnika słowem, nie samym atrybutem i nie samym krojem: przy czterech panelach
 * obok siebie sama kursywa jest nie do przeczytania bez porównania z sąsiadem.
 */
function listaStanow(
  uczestnik: RoundtableParticipant,
  indeks: number,
  powtorzony: boolean,
  cisza: boolean,
): HTMLElement {
  const stany = document.createElement('ul');
  stany.className = 'dr-uczestnik__stany';
  const opisy = [
    `Miejsce #${indeks + 1} w składzie znanym oknu`,
    cisza ? 'Wyciszony w turze — pytanie do niego nie idzie' : 'Słyszany w turze',
  ];
  if (powtorzony) opisy.push('Ten sam kanał modelu pod inną tożsamością — odrębny głos');
  if (uczestnik.order !== undefined) opisy.push(`Kolejność głosu wg rdzenia: #${uczestnik.order}`);
  // Pole klucza ustawia wyłącznie rdzeń — Operator nie ma czym go przestawić, panel je tylko pokazuje.
  if (uczestnik.key === true) opisy.push('Oznaczony przez rdzeń jako kluczowy');
  for (const opis of opisy) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-uczestnik__stan';
    pozycja.textContent = opis;
    stany.append(pozycja);
  }
  return stany;
}

/**
 * Prompt systemowy tożsamości — pokazywany, bo kontrakt go niesie i jest jedynym miejscem
 * odróżniającym dwie tożsamości jednego kanału.
 */
function promptSystemowy(uczestnik: RoundtableParticipant): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dr-uczestnik__prompt';
  const prompt = uczestnik.systemPrompt ?? '';
  element.dataset['prompt'] = prompt === '' ? 'brak' : 'jest';
  element.textContent =
    prompt === ''
      ? 'Prompt systemowy tożsamości: nie nadany — uczestnik odpowiada wedle nastaw swojego kanału.'
      : `Prompt systemowy tożsamości: ${prompt}`;
  return element;
}

/**
 * Metryki odpowiedzi uczestnika — długość i czas zmierzone z kontraktu, wraz ze zdaniem
 * o metrykach niemierzalnych.
 */
function metrykiPanelu(
  tekst: string,
  wypowiedz: RoundtableStatement | null,
  stan: StanDebaty,
): HTMLElement {
  const wykaz = document.createElement('ul');
  wykaz.className = 'dr-uczestnik__stany';
  wykaz.setAttribute('aria-label', 'Metryki odpowiedzi uczestnika');
  for (const zdanie of zdaniaMetryk(zmierzOdpowiedz(tekst, wypowiedz, stan.definicjaTury()))) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-uczestnik__stan';
    pozycja.textContent = zdanie;
    wykaz.append(pozycja);
  }
  return wykaz;
}

/**
 * Odpowiedź uczestnika — treść rosnąca albo utrwalona, zawsze ze zdaniem o tym, którą z nich
 * Operator ogląda.
 */
function odpowiedzUczestnika(biezacy: GlosBiezacy, cisza: boolean): HTMLElement {
  const blok = document.createElement('div');
  blok.className = 'dr-odpowiedz';
  blok.dataset['stanGlosu'] = biezacy.stan;
  blok.dataset['rosnie'] = biezacy.rosnie ? 'tak' : 'nie';

  const stanGlosu = document.createElement('p');
  stanGlosu.className = 'dr-odpowiedz__stan';
  stanGlosu.textContent = biezacy.zeStrumienia
    ? `${biezacy.stan} — treść rośnie ze strumienia, rdzeń jeszcze jej nie utrwalił`
    : biezacy.stan;
  blok.append(stanGlosu);

  if (biezacy.tekst === '') {
    const uwaga = document.createElement('p');
    uwaga.className = 'dr-odpowiedz__uwaga';
    uwaga.setAttribute('role', 'note');
    uwaga.textContent = zdanieBezSlow(biezacy, cisza);
    blok.append(uwaga);
  } else {
    const slowa = document.createElement('p');
    slowa.className = 'dr-odpowiedz__tresc';
    slowa.textContent = biezacy.tekst;
    blok.append(slowa);
  }

  if (biezacy.rozumowanie !== '') {
    const rozumowanie = document.createElement('p');
    rozumowanie.className = 'dr-odpowiedz__rozumowanie';
    rozumowanie.textContent = biezacy.rozumowanie;
    blok.append(rozumowanie);
  }
  if (biezacy.przyczyna !== '') {
    const zerwanie = document.createElement('p');
    zerwanie.className = 'dr-odpowiedz__uwaga';
    zerwanie.setAttribute('role', 'note');
    zerwanie.textContent = `Strumień zerwany: ${biezacy.przyczyna}`;
    blok.append(zerwanie);
  }
  return blok;
}
