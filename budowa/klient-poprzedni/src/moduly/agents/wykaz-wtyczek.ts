import type { Agent, AgentPlugin } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanAgentow } from './stan-agentow';

/**
 * Wykaz wtyczek eksperta — byt odrębny od konektora.
 *
 * Konektor jest drogą do usługi: serwerem MCP wskazanym punktem dostępu,
 * podawanym powłoce przez `--mcp-config`. Wtyczka jest katalogiem rozszerzeń
 * powłoki: zbiorem poleceń, zaczepów i umiejętności wgrywanym przez
 * `--plugin-dir`. Kontrakt rozdziela oba byty (`AgentPlugin`,
 * `Agent.pluginIds`, `agent.plugin.add`, `agent.plugin.remove`), więc
 * rozdziela je i okno.
 *
 * `Agent.pluginIds` niesie same identyfikatory, dlatego nazwę, wersję i źródło
 * wiersz bierze z `agent.plugin.list`. Identyfikator zostaje w podpowiedzi
 * wiersza, bo to nim posługuje się czynność odłączenia. Gdy odczyt definicji
 * odmówi albo jeszcze nie wrócił, wiersze pokazują same identyfikatory —
 * uboższa treść nie jest awarią i nie odbiera przycisku odłączenia.
 */
export interface WykazWtyczek {
  element: HTMLElement;
  /** Nanosi wtyczki eksperta czynnego; `null` znaczy brak eksperta. */
  ustaw(ekspert: Agent | null): void;
}

export function utworzWykazWtyczek(stan: StanAgentow): WykazWtyczek {
  const lista = document.createElement('ul');
  lista.className = 'da-wtyczki';

  const nazwa = poleTekstowe({ etykieta: 'Nazwa wtyczki', podpowiedz: 'np. danaco-narzedzia' });
  const zrodlo = poleTekstowe({
    etykieta: 'Katalog albo źródło wtyczki',
    podpowiedz: 'np. /opt/danaco/pluginy/narzedzia',
    opis: 'Katalog podawany powłoce przełącznikiem --plugin-dir.',
  });
  const wersja = poleTekstowe({ etykieta: 'Wersja wtyczki', podpowiedz: 'np. 1.4.0' });

  const podlacz = przycisk('Podłącz wtyczkę', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Wtyczki — katalog rozszerzeń powłoki';

  const nota = document.createElement('p');
  nota.className = 'dn-pole-opis';
  nota.textContent =
    'Wtyczka to nie konektor. Konektor jest drogą do usługi (--mcp-config), wtyczka ' +
    'katalogiem rozszerzeń powłoki (--plugin-dir). Ekspert może mieć jedno bez drugiego.';

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis da-granica';
  granica.textContent =
    'Nazwa, wersja i źródło pochodzą z agent.plugin.list. Gdy rdzeń nie odpowie, ' +
    'wiersz pokazuje sam identyfikator.';

  const element = document.createElement('section');
  element.className = 'da-panel da-wtyczki-panel';
  element.append(
    tytul,
    nota,
    lista,
    nazwa.element,
    zrodlo.element,
    wersja.element,
    podlacz,
    odpowiedz.element,
    granica,
  );

  let ekspertCzynny: Agent | null = null;
  /** Definicje wtyczek eksperta pokazanego, klucz = identyfikator wtyczki. */
  let definicje = new Map<string, AgentPlugin>();

  async function podlaczenie(): Promise<void> {
    if (ekspertCzynny === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — wtyczka należy do jednego.', false);
      return;
    }
    if (nazwa.kontrolka.value.trim() === '') {
      odpowiedz.pokaz('Wtyczka wymaga nazwy — rdzeń odmówi podłączenia bez niej.', false);
      return;
    }
    odpowiedz.pokaz('Podłączanie wtyczki w toku…', true);
    const wynik = await stan.zrodlo.dodajWtyczke({
      idEksperta: ekspertCzynny.id,
      nazwa: nazwa.kontrolka.value,
      zrodlo: zrodlo.kontrolka.value,
      wersja: wersja.kontrolka.value,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Podłączenie wtyczki', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const wtyczka = wynik.wynik.plugin;
    // Formularz czyścimy w całości: źródło poprzedniej wtyczki czekałoby w polu
    // i poszłoby do rdzenia przy następnym podłączeniu, także przy innym
    // ekspercie.
    nazwa.kontrolka.value = '';
    zrodlo.kontrolka.value = '';
    wersja.kontrolka.value = '';
    await stan.odswiez();
    odpowiedz.pokaz(
      `Wtyczka „${wtyczka.name}"${wtyczka.version === undefined ? '' : ` ${wtyczka.version}`} ` +
        `podłączona pod identyfikatorem ${wtyczka.id}.`,
      true,
    );
  }

  async function odlaczenie(idWtyczki: string): Promise<void> {
    if (ekspertCzynny === null) return;
    odpowiedz.pokaz('Odłączanie wtyczki w toku…', true);
    const wynik = await stan.zrodlo.odlaczWtyczke(ekspertCzynny.id, idWtyczki);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Odłączenie wtyczki', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    stan.wchlon(wynik.wynik.agent);
    odpowiedz.pokaz(`Wtyczka ${idWtyczki} odłączona od eksperta.`, true);
  }

  podlacz.addEventListener('click', () => void podlaczenie());

  /** Ekspert, którego wtyczki stoją w wykazie — po nim poznajemy zmianę wyboru. */
  let pokazany = '';

  return {
    element,

    ustaw(ekspert) {
      ekspertCzynny = ekspert;
      if ((ekspert?.id ?? '') !== pokazany) {
        pokazany = ekspert?.id ?? '';
        nazwa.kontrolka.value = '';
        zrodlo.kontrolka.value = '';
        wersja.kontrolka.value = '';
        odpowiedz.wyczysc();
      }
      const kody = ekspert?.pluginIds ?? [];
      if (kody.length > 0 && ekspert !== null) void wczytajDefinicje(ekspert.id);
      if (kody.length === 0) {
        definicje = new Map();
        lista.replaceChildren(
          pusty(
            ekspert === null
              ? 'Brak eksperta czynnego — wtyczka nie ma właściciela.'
              : `Ekspert „${ekspert.name}" nie ma jeszcze podłączonej wtyczki.`,
          ),
        );
        return;
      }
      lista.replaceChildren(
        ...kody.map((kod) => wiersz(kod, definicje.get(kod), () => void odlaczenie(kod))),
      );
    },
  };

  /**
   * Dobiera definicje wtyczek eksperta i przerysowuje wykaz.
   *
   * Odpowiedź przedawniona nie nadpisuje świeższej: między wysłaniem a powrotem
   * Operator może wybrać innego eksperta, a wtedy wykaz pokazywałby cudze
   * wtyczki pod właściwymi kodami.
   */
  async function wczytajDefinicje(idEksperta: string): Promise<void> {
    const wynik = await stan.zrodlo.wtyczki(idEksperta);
    if (!wynik.udany || wynik.wynik === undefined) return;
    if ((ekspertCzynny?.id ?? '') !== idEksperta) return;
    definicje = new Map(wynik.wynik.plugins.map((wtyczka) => [wtyczka.id, wtyczka]));
    const kody = ekspertCzynny?.pluginIds ?? [];
    lista.replaceChildren(
      ...kody.map((kod) => wiersz(kod, definicje.get(kod), () => void odlaczenie(kod))),
    );
  }
}

function wiersz(
  kod: string,
  definicja: AgentPlugin | undefined,
  naOdlaczenie: () => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'da-wtyczki__wiersz';
  element.dataset['wtyczka'] = kod;

  const etykieta = document.createElement('span');
  etykieta.className = 'da-wtyczki__kod';
  // Bez definicji zostaje sam identyfikator. Napis zastępczy w rodzaju
  // „(bez nazwy)" udawałby, że rdzeń odpowiedział.
  etykieta.textContent =
    definicja === undefined
      ? kod
      : `${definicja.name}${definicja.version === undefined ? '' : ` ${definicja.version}`}`;
  if (definicja !== undefined) {
    etykieta.title = definicja.source === undefined ? kod : `${kod} · ${definicja.source}`;
  }

  const odlacz = przycisk('Odłącz', 'dn-btn dn-btn--sm dn-btn--zarys');
  odlacz.addEventListener('click', naOdlaczenie);

  element.append(etykieta, odlacz);
  return element;
}

function pusty(tresc: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'da-wtyczki__wiersz';
  element.dataset['pusty'] = 'tak';
  element.textContent = tresc;
  return element;
}
