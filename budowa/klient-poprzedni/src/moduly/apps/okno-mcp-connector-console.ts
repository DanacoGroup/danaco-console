import { ExtensionKind } from '../../../../shared/contract';
import { poleTekstowe, poleWielowierszowe } from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import {
  BRAKI_MCP_CONSOLE,
  KODY_OKIEN,
  KOMENDY_PROPONOWANE,
  NAZWY_OKIEN,
  OKNO_SPOZA_KATALOGU,
} from './etykiety-apps';
import { utworzRameApps } from './rama-okna';
import type { StanRozszerzen } from './stan-rozszerzen';

/**
 * MCP & Connector Console — panel boczny inspektora narzędzi i diagnostyki
 * protokołu.
 *
 * Okno stoi w całości na czynnościach, których kontrakt nie prowadzi, i jest
 * zbudowane właśnie po to, żeby ten brak był widoczny i policzony. Kontrolki są
 * obecne, klikalne i nazywają brakującą komendę; żadna nie jest wygaszona
 * i żadna nie udaje wykonania.
 *
 * Powód każdej kontrolki bierze się z bytu pokrycia (`moduly/pokrycie-komend.ts`),
 * a nie z napisu wpisanego tutaj: byt pyta rdzeń o wykaz jego komend i sam
 * przerysowuje zdanie, więc w dniu, w którym rdzeń dostanie uchwyt dla którejś
 * z tych komend, kontrolka powie to sama, bez tknięcia tego pliku. To ważne, bo
 * brak jest stanem przejściowym — komendy są zaprojektowane do natychmiastowej
 * implementacji.
 *
 * Formularz argumentów jest jednym polem tekstowym, nie polami wyprowadzonymi
 * ze schematu: schemat przychodzi z odkrycia narzędzi, którego nie ma. Pola
 * zmyślone ze zgadniętego schematu wyglądałyby na wiedzę o serwerze, której
 * okno nie ma.
 *
 * Okno otwiera się na pozycji wskazanej w Integrations Hubie. Inspektor bez
 * wskazanej integracji nie ma czego inspekcjonować i mówi to wprost.
 */
export interface OknoMcpConnectorConsole {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoMcpConnectorConsole(
  stan: StanRozszerzen,
  pokrycie: PokrycieKomend,
): OknoMcpConnectorConsole {
  const kod = KODY_OKIEN.McpConnectorConsole;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'pomocnicze');

  const wskazana = document.createElement('p');
  wskazana.className = 'mp-konsola__wskazana';

  const narzedzia = pokrycie.przycisk(
    'Odkryj narzędzia, zasoby i prompty',
    KOMENDY_PROPONOWANE.OdkrycieNarzedzi,
    'Pobranie wykazu tools, resources i prompts udostępnianych przez serwer wraz ze schematami',
  );

  const nazwaNarzedzia = poleTekstowe({
    etykieta: 'Nazwa narzędzia',
    podpowiedz: 'nazwa z wykazu serwera',
  });
  const argumenty = poleWielowierszowe(
    { etykieta: 'Argumenty wywołania (JSON)', podpowiedz: '{"query":"…"}' },
    4,
  );

  const wywolaj = pokrycie.przycisk(
    'Wywołaj próbnie',
    KOMENDY_PROPONOWANE.WywolanieProbne,
    'Wykonanie tools/call na wskazanym narzędziu z podglądem odpowiedzi surowej i sformatowanej',
  );
  const piaskownica = pokrycie.przycisk(
    'Uruchom w piaskownicy',
    KOMENDY_PROPONOWANE.PiaskownicaTestu,
    'Uruchomienie umiejętności albo wtyczki na przykładowym wejściu bez podłączania jej do eksperta',
  );
  const dziennik = pokrycie.przycisk(
    'Pokaż log JSON-RPC',
    KOMENDY_PROPONOWANE.DziennikProtokolu,
    'Strumień żądań, odpowiedzi i notyfikacji wymienianych z serwerem',
  );

  const wykazPokrycia = pokrycie.wykaz(
    [
      KOMENDY_PROPONOWANE.OdkrycieNarzedzi,
      KOMENDY_PROPONOWANE.WywolanieProbne,
      KOMENDY_PROPONOWANE.PiaskownicaTestu,
      KOMENDY_PROPONOWANE.DziennikProtokolu,
    ],
    'mp-granica mp-pokrycie',
  );

  const pasekNarzedzi = document.createElement('div');
  pasekNarzedzi.className = 'mp-pasek';
  pasekNarzedzi.append(narzedzia, dziennik);

  const pasekWywolania = document.createElement('div');
  pasekWywolania.className = 'mp-pasek';
  pasekWywolania.append(wywolaj, piaskownica);

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis mp-granica';
  granica.textContent = OKNO_SPOZA_KATALOGU;

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_MCP_CONSOLE, 'extension'));
  rama.tresc.append(
    wskazana,
    naglowekCzesci('Odkryte definicje serwera'),
    pasekNarzedzi,
    naglowekCzesci('Próbne wywołanie'),
    opiszPole(
      nazwaNarzedzia.element,
      'Nazwa wpisywana ręcznie, bo wykazu narzędzi serwera nie ma skąd wziąć — komendy ' +
        'odkrycia kontrakt nie niesie.',
    ),
    opiszPole(
      argumenty.element,
      'Jedno pole tekstowe zamiast pól wyprowadzonych ze schematu JSON Schema: schemat ' +
        'przychodzi z odkrycia narzędzi, którego nie ma. Pola zgadnięte udawałyby wiedzę ' +
        'o serwerze, której okno nie ma.',
    ),
    pasekWywolania,
    wykazPokrycia,
    granica,
  );

  function odswiez(): void {
    const pozycja = stan.wybrane();
    if (pozycja === null) {
      wskazana.textContent =
        'Nie wskazano integracji. Naciśnij „Konsola i inspektor" w wierszu Integrations Hubu.';
      if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
      rama.puste(
        'Inspektor bez wskazanej integracji nie ma czego inspekcjonować. Wskaż serwer MCP ' +
          'albo integrację API w Integrations Hubie.',
      );
      return;
    }
    // Konsola dotyczy serwera MCP; dla pozostałych rodzajów mówi, czego dotyczy
    // i czego nie — zamiast udawać, że wtyczka odpowiada na tools/list.
    const oRodzaju =
      pozycja.kind === ExtensionKind.Mcp
        ? 'Log JSON-RPC i odkrycie definicji dotyczą tego rodzaju pozycji wprost.'
        : `Pozycja jest rodzaju ${pozycja.kind}: protokołu JSON-RPC nie prowadzi, więc log ` +
          'protokołu i odkrycie narzędzi jej nie dotyczą. Piaskownica testu dotyczy.';
    wskazana.textContent =
      `Wskazana integracja: ${pozycja.name} (kod ${pozycja.code}, rodzaj ${pozycja.kind}, ` +
      `${pozycja.enabled ? 'włączona' : 'wyłączona'}). ${oRodzaju}`;
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    // Okno nie przechodzi w stan „gotowe”, bo niczego nie odczytało: stan pusty
    // z nazwanym powodem jest tu stanem prawdziwym.
    rama.puste(
      'Konsola nie ma dziś ani jednej drogi do rdzenia. Cztery czynności — odkrycie ' +
        'definicji, próbne wywołanie, piaskownica i log protokołu — czekają na komendy; ' +
        'przyciski nazywają je po naciśnięciu.',
    );
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Nagłówek części okna. */
function naglowekCzesci(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mp-czesc__tytul';
  element.textContent = tresc;
  return element;
}
