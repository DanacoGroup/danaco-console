import { IdentityLayer, IdentityMode, type Agent, type AgentLayer } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleLogiczne,
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanAgentow } from './stan-agentow';

/**
 * Edytor trzech warstw promptu eksperta — konstytucja, profil roli, ekspertyza.
 *
 * Warstwy są uporządkowane wg krytyczności: konstytucja stoi najwyżej,
 * ekspertyza najniżej. Kolejność rozstrzyga, co wygrywa przy sprzeczności
 * między warstwami, więc stoi przy każdej warstwie liczbą i zdaniem, a nie
 * w domyśle wynikającym z porządku pól na ekranie. Nakładkę składa rdzeń
 * (`server/internal/injection/nakladka.go`), a kontrakt niesie `AgentLayer`
 * i cztery komendy `agent.layer.*`.
 *
 * Warstwa nie ma własnego trybu, bo dopisanie do globalnego promptu systemowego
 * albo jego zastąpienie dotyczy instrukcji eksperta jako całości: przełącznik
 * stoi raz, przy tożsamości (`formularz-tozsamosci.ts`). Dlatego `AgentLayer`
 * i `agent.layer.set` pola `mode` nie mają, a `Agent.mode` ma. W miejscu
 * przełącznika stoi przy każdej warstwie zdanie czytające tryb eksperta czynnego.
 *
 * Edytor nie dotyka `Agent.systemPrompt` — to pole płaskie zapisuje formularz
 * tożsamości komendą `agent.update`.
 */
export interface WarstwyPromptu {
  element: HTMLElement;
  /** Nanosi warstwy eksperta czynnego; `null` znaczy brak eksperta. */
  ustaw(ekspert: Agent | null): void;
}

/** Kolejność warstw wg krytyczności — konstytucja najwyżej. */
export const KOLEJNOSC_WARSTW: readonly IdentityLayer[] = [
  IdentityLayer.Constitution,
  IdentityLayer.Profile,
  IdentityLayer.Expertise,
];

/** Nazwy warstw w języku Operatora. */
export const NAZWY_WARSTW: Record<IdentityLayer, string> = {
  [IdentityLayer.Constitution]: 'konstytucja',
  [IdentityLayer.Profile]: 'profil roli',
  [IdentityLayer.Expertise]: 'ekspertyza zadaniowa',
};

const RANGI_WARSTW: Record<IdentityLayer, string> = {
  [IdentityLayer.Constitution]: 'warstwa najwyższa — wygrywa z pozostałymi',
  [IdentityLayer.Profile]: 'warstwa środkowa — ustępuje konstytucji',
  [IdentityLayer.Expertise]: 'warstwa najniższa — ustępuje obu wyższym',
};

/**
 * Czy ekspert ma oznaczone odstępstwo od ustawień domyślnych.
 *
 * `Agent.mode` pominięte znaczy `DOLACZ` — stan domyślny. Ekspert jeszcze
 * niezałożony (`null`) czyta się tak samo: `agent.create` bez pola `mode`
 * zakłada eksperta dopisującego.
 */
export function czyZastepuje(ekspert: Agent | null): boolean {
  return ekspert?.mode === IdentityMode.ZASTAP;
}

/** Zdanie przy panelu: stan domyślny oraz droga odstąpienia od niego. */
const DROGA_DO_MODELU =
  'Domyślnie warstwy eksperta DOPISUJĄ się do globalnego promptu systemowego ' +
  'z okna konfiguracji i ustawień na stronie głównej — globalny obowiązuje pierwszy. ' +
  'Od tego stanu Operator może odstąpić przy tożsamości eksperta wyżej: po oznaczeniu ' +
  'odstępstwa instrukcja tego eksperta STAJE SIĘ promptem systemowym, a globalny ' +
  'przestaje obowiązywać dla jego okien.';

/**
 * Zdanie przy warstwie — zależne od trybu eksperta czynnego. Zdanie stałe
 * byłoby nieprawdziwe przy jednym z dwóch stanów.
 */
const DROGA_WARSTWY: Record<'dopisanie' | 'zastapienie', string> = {
  dopisanie:
    'Ta warstwa dopisuje się do globalnego promptu systemowego — globalny obowiązuje ' +
    'pierwszy. Odstępstwo od tego oznacza się raz, przy tożsamości eksperta.',
  zastapienie:
    'Ekspert ma oznaczone odstępstwo od ustawień domyślnych: ta warstwa wchodzi do ' +
    'promptu systemowego ZAMIAST globalnego, który dla jego okien nie obowiązuje.',
};

export function utworzWarstwyPromptu(stan: StanAgentow): WarstwyPromptu {
  const odpowiedz = utworzWierszOdpowiedzi();
  const czesci = KOLEJNOSC_WARSTW.map((warstwa, numer) =>
    utworzWarstwe(stan, warstwa, numer + 1, odpowiedz.pokaz),
  );

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Warstwy promptu — kolejność krytyczności';

  const nota = document.createElement('p');
  nota.className = 'dn-pole-opis';
  nota.textContent =
    'Warstwy schodzą od najwyższej do najniższej. Przy sprzeczności treści obowiązuje ' +
    'warstwa wyższa; warstwa nieczynna nie wchodzi do promptu wcale.';

  const droga = document.createElement('p');
  droga.className = 'dn-pole-opis da-warstwy-panel__droga';
  droga.textContent = DROGA_DO_MODELU;
  droga.setAttribute('role', 'note');

  const element = document.createElement('section');
  element.className = 'da-panel da-warstwy-panel';
  element.append(tytul, nota, droga, ...czesci.map((czesc) => czesc.element), odpowiedz.element);

  return {
    element,

    ustaw(ekspert) {
      for (const czesc of czesci) czesc.ustaw(ekspert);
      if (ekspert === null) odpowiedz.wyczysc();
    },
  };
}

/** Jedna warstwa: treść, czynność oraz zapis i usunięcie. */
interface CzescWarstwy {
  element: HTMLElement;
  ustaw(ekspert: Agent | null): void;
}

function utworzWarstwe(
  stan: StanAgentow,
  warstwa: IdentityLayer,
  numer: number,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): CzescWarstwy {
  const tresc = poleWielowierszowe(
    { etykieta: `Treść warstwy — ${NAZWY_WARSTW[warstwa]}`, podpowiedz: 'treść wchodząca w prompt' },
    5,
  );
  const czynna = poleLogiczne({ etykieta: 'Warstwa czynna' });
  czynna.kontrolka.checked = true;

  const droga = document.createElement('p');
  droga.className = 'dn-pole-opis da-warstwa__droga';

  const zapisz = przycisk('Zapisz warstwę', 'dn-btn dn-btn--sm dn-btn--atrament');
  const usun = przycisk('Usuń treść warstwy', 'dn-btn dn-btn--sm dn-btn--zarys');

  const dzialania = document.createElement('div');
  dzialania.className = 'da-warstwa__dzialania';
  dzialania.append(zapisz, usun);

  const ranga = document.createElement('p');
  ranga.className = 'dn-pole-opis da-warstwa__ranga';
  ranga.textContent = `${numer}. ${NAZWY_WARSTW[warstwa]} — ${RANGI_WARSTW[warstwa]}`;

  const stanZapisu = document.createElement('p');
  stanZapisu.className = 'dn-pole-opis da-warstwa__stan';

  const element = document.createElement('div');
  element.className = 'da-warstwa';
  element.dataset['warstwa'] = warstwa;
  element.append(ranga, stanZapisu, tresc.element, droga, czynna.element, dzialania);

  let ekspertCzynny: Agent | null = null;

  async function zapisanie(): Promise<void> {
    if (ekspertCzynny === null) {
      powiedz('Warstwa należy do eksperta — załóż go najpierw w formularzu wyżej.', false);
      return;
    }
    powiedz(`Zapis warstwy „${NAZWY_WARSTW[warstwa]}" w toku…`, true);
    const wynik = await stan.zrodlo.zapiszWarstwe({
      idEksperta: ekspertCzynny.id,
      warstwa,
      tresc: tresc.kontrolka.value,
      czynna: czynna.kontrolka.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowy('Zapis warstwy promptu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    stan.wchlon(wynik.wynik.agent);
    powiedz(
      `Warstwa „${NAZWY_WARSTW[warstwa]}" zapisana — ${
        czyZastepuje(wynik.wynik.agent)
          ? 'wchodzi do promptu systemowego zamiast globalnego (odstępstwo oznaczone).'
          : 'dopisuje się do globalnego promptu systemowego.'
      }`,
      true,
    );
  }

  async function usuniecie(): Promise<void> {
    if (ekspertCzynny === null) {
      powiedz('Warstwa należy do eksperta — nie ma czego usuwać bez niego.', false);
      return;
    }
    powiedz(`Usuwanie warstwy „${NAZWY_WARSTW[warstwa]}" w toku…`, true);
    const wynik = await stan.zrodlo.usunWarstwe(ekspertCzynny.id, warstwa);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowy('Usunięcie warstwy promptu', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    stan.wchlon(wynik.wynik.agent);
    powiedz(`Warstwa „${NAZWY_WARSTW[warstwa]}" usunięta — prompt idzie bez niej.`, true);
  }

  zapisz.addEventListener('click', () => void zapisanie());
  usun.addEventListener('click', () => void usuniecie());

  return {
    element,

    ustaw(ekspert) {
      ekspertCzynny = ekspert;
      const zapisana: AgentLayer | undefined = ekspert?.layers?.find(
        (wpis) => wpis.layer === warstwa,
      );
      tresc.kontrolka.value = zapisana?.content ?? '';
      czynna.kontrolka.checked = zapisana?.enabled ?? true;
      const zastepuje = czyZastepuje(ekspert);
      droga.textContent = DROGA_WARSTWY[zastepuje ? 'zastapienie' : 'dopisanie'];
      // Odstępstwo ma wagę wizualną także tutaj — ten sam znacznik `data-`,
      // którym rządzi arkusz sekcji modeli (`modele/tozsamosc.css`).
      droga.dataset['zastapienie'] = String(zastepuje);
      stanZapisu.textContent =
        ekspert === null
          ? 'Brak eksperta czynnego — warstwa nie ma właściciela.'
          : zapisana === undefined
            ? 'Warstwa niezapisana — ekspert idzie dziś bez niej.'
            : `Warstwa zapisana ${new Date(zapisana.updatedAt).toLocaleString('pl-PL')}.`;
    },
  };
}
