import { Command, ProviderTransport, type Channel } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWielowierszowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';
import { naglowek } from './biblioteka-ekspertow';
import { utworzPodgladWywolania } from './podglad-wywolania';
import { utworzPolaZalezne, type PolaZalezne } from './pola-zalezne';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAgentow } from './stan-agentow';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Model Configuration to okno pomocnicze modułu Agents — ustawia kanał, transport, model
 * bazowy i parametry wywołania eksperta.
 */
export interface OknoModelConfiguration {
  element: HTMLElement;
  odswiez(): void;
  /** Odczytuje rejestr kanałów i deklarację zdolności adaptera. */
  wczytaj(): Promise<void>;
  /** Odpina wykaz pokrycia od wspólnego bytu. Wołane z `rozlacz()` modułu. */
  zamknij(): void;
}

/**
 * Formatuje parametry wywołania eksperta do postaci tekstowej wpisywanej w polu
 * formularza, z pustym wynikiem dla wartości nieustawionej.
 */
function zapisParametrow(wartosc: unknown): string {
  if (wartosc === undefined || wartosc === null) return '';
  return JSON.stringify(wartosc, null, 2);
}

export function utworzOknoModelConfiguration(
  stan: StanAgentow,
  zaplecze: ZrodloZaplecza,
  kanalRdzenia: Kanal,
): OknoModelConfiguration {
  const okno: StanOkna = utworzStanOkna();
  const zdolnosci: PolaZalezne = utworzPolaZalezne();
  const pokrycie: PokrycieKomend = utworzPokrycieKomend(kanalRdzenia);

  const kanal = poleWyboru({ etykieta: 'Kanał modelu (rejestr rdzenia)' }, []);
  const transport = poleWyboru(
    {
      etykieta: 'Droga wywołania kanału',
      opis: 'Puste znaczy „ekspert nie ma własnej drogi” — obowiązuje wtedy transport kanału.',
    },
    [
      { wartosc: '', etykieta: 'transport kanału' },
      ...Object.values(ProviderTransport).map((wartosc) => ({ wartosc, etykieta: wartosc })),
    ],
  );
  const model = poleTekstowe({
    etykieta: 'Model bazowy',
    opis: 'Puste znaczy model wskazany przez wiersz kanału.',
  });
  const parametry = poleWielowierszowe(
    {
      etykieta: 'Parametry wywołania (JSON)',
      podpowiedz: '{"maxOutputTokens": 8000}',
      opis: 'Parametry zapisane przy tym ekspercie; puste znaczy wywołanie bez nich.',
    },
    4,
  );

  const zapisz = przycisk('Zapisz model bazowy', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  // Test połączenia sprawdza most, nie kanał modelu — kontrolka zostaje widoczna i nazywa brak.
  const test = pokrycie.przycisk(
    'Testuj połączenie',
    Command.ChannelCheck,
    'sprawdzenie osiągalności kanału modelu przed zapisem eksperta',
  );

  // Poświadczenie kanału zakłada się i zmienia w oknie konfiguracji platformy, nie tutaj.
  const poswiadczenie = pokrycie.przycisk(
    'Stan poświadczenia kanału',
    Command.ChannelCredentialStatus,
    'odczyt stanu poświadczenia kanału — ustawione czy brak i kiedy zmienione, nigdy treść klucza',
  );

  const perRola = document.createElement('p');
  perRola.className = 'dn-pole-opis da-granica';
  perRola.textContent =
    'Wartość zapisana tutaj jest domyślną wartością eksperta. Przy przypisaniu eksperta ' +
    'do roli w środowisku MultitaskingAI kanał modelu może zostać dla tej roli nadpisany — ' +
    'poza rolą obowiązuje ustawienie z tego okna.';

  const formularz = document.createElement('div');
  formularz.className = 'da-model__formularz';
  formularz.append(
    kanal.element,
    transport.element,
    model.element,
    parametry.element,
    zapisz,
    test,
    poswiadczenie,
  );
  // Podgląd wywołania pokazuje wiersz polecenia i prompt systemowy po nałożeniu eksperta.
  const podglad = utworzPodgladWywolania(kanalRdzenia);
  okno.tresc.append(formularz, zdolnosci.element, odpowiedz.element, perRola, podglad.element);

  const element = document.createElement('section');
  element.className = 'da-okno da-okno--pomocnicze';
  element.dataset['okno'] = 'model-configuration';
  element.append(naglowek('Model Configuration'), okno.element);

  let rejestr: Channel[] = [];
  /** Powód odmowy odczytu rejestru kanałów; pusty, gdy odczyt się udał. */
  let odmowaRejestru = '';
  /** Ekspert, którego wartości stoją w formularzu — po nim poznajemy zmianę. */
  let pokazany = '';

  async function wczytajZdolnosci(): Promise<void> {
    zdolnosci.ladowanie();
    const wynik = await zaplecze.zdolnosci(kanal.kontrolka.value, wybranyTransport());
    if (!wynik.udany || wynik.wynik === undefined) {
      zdolnosci.blad(
        opisOdmowy('Odczyt zdolności adaptera', wynik.blad?.code, wynik.blad?.message),
      );
      return;
    }
    zdolnosci.ustaw(wynik.wynik.capabilities);
  }

  function wybranyTransport(): ProviderTransport | '' {
    return transport.kontrolka.value as ProviderTransport | '';
  }

  async function zapisanie(): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — konfiguracja dotyczy jednego.', false);
      return;
    }
    if (kanal.kontrolka.value === '') {
      odpowiedz.pokaz('Rdzeń wymaga wskazania kanału modelu.', false);
      return;
    }
    odpowiedz.pokaz('Zapis modelu bazowego w toku…', true);
    const wynik = await stan.zrodlo.ustawModel({
      idEksperta: ekspert.id,
      kanal: kanal.kontrolka.value,
      model: model.kontrolka.value,
      transport: wybranyTransport(),
      parametry: parametry.kontrolka.value,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zapis modelu bazowego', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    // Transport i parametry nie wracają w polu Agent, więc zdanie powodzenia wymienia, co wysłano.
    const wyslane: string[] = [];
    if (wybranyTransport() !== '') wyslane.push(`transport ${wybranyTransport()}`);
    if (parametry.kontrolka.value.trim() !== '') wyslane.push('parametry wywołania');

    // Formularz przyjmuje stan zapisanego eksperta z rdzenia, nie to, co wpisał operator.
    const zapisany = wynik.wynik.agent;
    stan.wchlon(zapisany);
    pokazany = zapisany.id;
    transport.kontrolka.value = zapisany.transport ?? '';
    parametry.kontrolka.value = zapisParametrow(zapisany.parameters);
    odpowiedz.pokaz(
      `Model bazowy zapisany — kanał ${zapisany.channelId ?? '—'}, ` +
        `model ${zapisany.model ?? '—'}, droga ${zapisany.transport ?? 'kanału'}.` +
        (wyslane.length === 0 ? '' : ` Wysłano też ${wyslane.join(' i ')}.`),
      true,
    );
  }

  function odswiez(): void {
    const ekspert = stan.wybrany();
    if (ekspert === null) {
      // Odmowa odczytu rejestru zostaje na widoku zamiast podpowiedzi o wyborze eksperta.
      pokazany = '';
      if (odmowaRejestru !== '') okno.blad(odmowaRejestru);
      else okno.puste('Wybierz eksperta w Agent Builderze, aby ustalić jego model bazowy.');
      return;
    }
    if (ekspert.id !== pokazany) {
      // Zmiana eksperta przestawia formularz na stan rdzenia, żadnego pola nie zeruje w ciemno.
      pokazany = ekspert.id;
      transport.kontrolka.value = ekspert.transport ?? '';
      parametry.kontrolka.value = zapisParametrow(ekspert.parameters);
      odpowiedz.wyczysc();
    }
    if (rejestr.length === 0) return;
    kanal.kontrolka.value = ekspert.channelId ?? '';
    model.kontrolka.value = ekspert.model ?? '';
    okno.gotowe();
  }

  kanal.kontrolka.addEventListener('change', () => void wczytajZdolnosci());
  transport.kontrolka.addEventListener('change', () => void wczytajZdolnosci());
  zapisz.addEventListener('click', () => void zapisanie());

  return {
    element,

    odswiez,

    zamknij: () => pokrycie.zamknij(),

    async wczytaj() {
      // Odczyt wykazu komend jest zlecany swobodnie, niezależnie od powitania połączenia.
      void pokrycie.odczytaj();
      okno.ladowanie('Odczyt rejestru kanałów modelu w toku…');
      const wynik = await zaplecze.kanaly();
      if (!wynik.udany || wynik.wynik === undefined) {
        odmowaRejestru = opisOdmowy(
          'Odczyt rejestru kanałów',
          wynik.blad?.code,
          wynik.blad?.message,
        );
        okno.blad(odmowaRejestru);
        return;
      }
      odmowaRejestru = '';
      rejestr = wynik.wynik.channels;
      ustawPozycje(kanal.kontrolka, [
        { wartosc: '', etykieta: 'bez kanału bazowego' },
        ...rejestr.map((wiersz) => ({
          wartosc: wiersz.id,
          etykieta: `${wiersz.name} · ${wiersz.kind}${wiersz.model === undefined ? '' : ` · ${wiersz.model}`}`,
        })),
      ]);
      if (rejestr.length === 0) {
        okno.puste('Rejestr kanałów modelu jest pusty — załóż kanał w oknie konfiguracji.');
        return;
      }
      okno.gotowe();
      odswiez();
      await wczytajZdolnosci();
    },
  };
}
