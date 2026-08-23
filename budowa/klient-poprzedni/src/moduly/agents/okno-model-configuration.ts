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
 * Model Configuration — okno pomocnicze modułu Agents.
 *
 * Lista kanałów pochodzi z komendy `channel.list` — tego samego rejestru,
 * z którego biorą kanał okna rozmowy. Okno nie ma własnej listy dostawców
 * i nie zna nazw programów CLI; kanały „Code CLI”, „Agent SDK” i „API” są
 * wierszami tego rejestru wraz z transportem, nie gałęziami w kodzie.
 *
 * Zmiana kanału przełącza zestaw pól zależnych. Zestaw bierze się z deklaracji
 * zdolności adaptera (`config.capabilities.get`, obszar `model` modelu
 * konfiguracji sesji), więc przełącza się sam, gdy zmieni się kanał albo
 * transport.
 *
 * Komendy okna: `agent.model.set` (zapis) oraz `channel.list`
 * i `config.capabilities.get` (odczyt). Komenda `model.channel.set` nadaje
 * kanał modelu oknu komunikacji albo karcie sesji, a Model Configuration
 * dotyczy jednego eksperta i żadnego okna komunikacji nie prowadzi — nie ma
 * czym wskazać `windowId`, więc okno jej nie woła.
 *
 * Wszystkie cztery pola są odczytywalne: `Agent` niesie `channelId`, `model`,
 * `transport` i `parameters`, a bazą trzyma je migracja agentów. Formularz
 * pokazuje więc stan rdzenia, a nie własną pamięć — przy zmianie eksperta pola
 * przyjmują wartości nowego, a nie zostają po poprzedniku. Puste pole znaczy
 * „ekspert tego nie ma”, i tak je opisuje.
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
 * Zapis parametrów wywołania do pola tekstowego.
 *
 * Parametry są w kontrakcie wartością JSON dowolnego kształtu, więc do pola
 * idą tekstem sformatowanym — nie `String(obiekt)`, bo to dałoby napis
 * „[object Object]” zamiast treści, którą Operator ma poprawić. Brak wartości
 * daje pole puste, a nie napis „undefined”.
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

  // Test połączenia nie ma dziś drogi do rdzenia: sprawdzenie punktu dostępu
  // dotyczy mostu, nie kanału modelu, a rodzina kanałów sprawdzenia nie niesie.
  // Kontrolka zostaje widoczna i klikalna — po naciśnięciu nazywa brak.
  const test = pokrycie.przycisk(
    'Testuj połączenie',
    Command.ChannelCheck,
    'sprawdzenie osiągalności kanału modelu przed zapisem eksperta',
  );

  // Poświadczenie kanału stoi poza modułem: zakłada je i zmienia rodzina
  // kanałów w oknie konfiguracji platformy. Okno mówi, gdzie ta decyzja
  // zapada, zamiast stawiać pole, które nie miałoby dokąd pojechać.
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
  // Podgląd wywołania zamyka pytanie „co to okno naprawdę robi”: pokazuje
  // wiersz polecenia i prompt systemowy po nałożeniu eksperta, liczone tą samą
  // drogą, którą idzie tura (`config.explain.get`).
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
    // Wysłane, ale nieodczytywalne: rdzeń nie oddaje tych dwóch pól w `Agent`,
    // więc zdanie powodzenia wymienia, co poszło, i nie udaje odczytu.
    const wyslane: string[] = [];
    if (wybranyTransport() !== '') wyslane.push(`transport ${wybranyTransport()}`);
    if (parametry.kontrolka.value.trim() !== '') wyslane.push('parametry wywołania');

    // Zapisany ekspert wraca z rdzenia z kompletem czterech pól, więc formularz
    // przyjmuje jego stan zamiast zostawać przy tym, co Operator wpisał —
    // rozjazd między jednym a drugim jest odpowiedzią rdzenia, nie usterką okna.
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
      // Odmowa odczytu rejestru zostaje na widoku. Podpowiedź „wybierz eksperta”
      // nie ma prawa zająć miejsca powodu, który rdzeń podał.
      pokazany = '';
      if (odmowaRejestru !== '') okno.blad(odmowaRejestru);
      else okno.puste('Wybierz eksperta w Agent Builderze, aby ustalić jego model bazowy.');
      return;
    }
    if (ekspert.id !== pokazany) {
      // Zmiana eksperta przestawia formularz na stan rdzenia dla nowo wybranego.
      // Wszystkie cztery pola są odczytywalne z bytu eksperta, więc żadnego nie
      // trzeba zerować „w ciemno” — po zmianie widać to, co rdzeń ma zapisane.
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
      // Powitanie idzie raz na połączenie, więc odczyt wykazu komend można
      // zlecić swobodnie — od niego zależy zdanie kontrolek bez pokrycia.
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
