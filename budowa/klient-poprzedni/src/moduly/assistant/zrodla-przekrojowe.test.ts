import { describe, expect, it } from 'vitest';

import { Command, ListenMode } from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzZrodloAssistant } from './zrodlo-assistant';
import { utworzZrodloKontekstow } from './zrodlo-kontekstow';
import { utworzZrodloMowy } from './zrodlo-mowy';
import { utworzZrodloSchowka } from './zrodlo-schowka';
import { utworzZrodloZaplecza } from './zrodlo-zaplecza';

/** Sprawdzian dróg klienckich rodzin przekrojowych obsługiwanych przez moduł Assistant. */

/** Kanał próbny zapamiętuje komendy i oddaje odpowiedź pustą, mierząc warstwę kliencką, nie zachowanie rdzenia. */
function kanalProbny(): { kanal: Kanal; wyslane: string[] } {
  const wyslane: string[] = [];
  const kanal = {
    wyslij(komenda: string, _zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push(komenda);
      przyWyniku?.({ udany: true, wynik: {} });
      return `zadanie-${String(wyslane.length)}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Wywołuje każdą czynność czterech źródeł raz, dając pełny przelot rodzin komend przekrojowych tego modułu. */
async function przelotZrodel(kanal: Kanal): Promise<void> {
  const asystent = utworzZrodloAssistant(kanal);
  await asystent.polecenie({
    idOkna: 'okno-1',
    transkrypcja: 'podsumuj pocztę',
    profil: '',
    czytaj: false,
  });
  await asystent.zlecenia('okno-1');
  await asystent.dziennik('okno-1', '');
  await asystent.oznaczWpis('wpis-1', true, 'do przeglądu');

  // Dwie komendy mowy zastane mieszkają w zapleczu modułu, bo tam stały, zanim rodzina urosła.
  const zaplecze = utworzZrodloZaplecza(kanal);
  await zaplecze.silnikMowy();
  await zaplecze.transkrypcja({ audioRef: '/dane/nagrania-mowy/nagranie-1.wav' });

  const mowa = utworzZrodloMowy(kanal);
  await mowa.przeslijNagranie({ base64: 'AAA=', typTresci: 'audio/wav', idOkna: 'okno-1' });
  await mowa.pobierzNagranie('/dane/nagrania-mowy/nagranie-1.wav');
  await mowa.nastawaWybudzania();
  await mowa.zapiszWybudzanie({ fraza: 'Danaco', tryb: ListenMode.WakeWord });
  await mowa.uruchomNasluch('okno-1', 'karta-1');
  await mowa.zatrzymajNasluch('nasluch-1', 'okno-1');

  const konteksty = utworzZrodloKontekstow(kanal);
  await konteksty.konteksty('karta-1', false);
  await konteksty.zapiszKontekst({
    id: '',
    nazwa: 'projekt Atlas',
    opis: '',
    poziomy: [],
    wpisy: [],
    promptSystemowy: '',
    czynny: true,
  });
  await konteksty.uaktywnijKontekst('kontekst-1', 'karta-1');
  await konteksty.usunKontekst('kontekst-1');
  await konteksty.zasadyRetencji();
  await konteksty.zapiszZasadeRetencji({
    dniWygasania: 30,
    wrazliweDomyslnie: false,
    wzorceNigdyNieZapisywane: [],
    czynna: true,
  });
  await konteksty.zajetoscKontekstu('okno-1', 'karta-1');

  const schowek = utworzZrodloSchowka(kanal);
  await schowek.historia('', false);
  await schowek.dopisz({ tresc: 'ustalenia' });
  await schowek.przypnij('wpis-1', true);
  await schowek.usun('wpis-1', false);
  await schowek.skroty('');
  await schowek.zapiszSkrot({
    id: '',
    skrot: ';odmowa',
    tresc: 'Dziękuję za propozycję.',
    opis: '',
    polaSzablonu: [],
    czynny: true,
  });
  await schowek.usunSkrot('skrot-1');
  await schowek.skrotGlobalny();
  await schowek.zapiszSkrotGlobalny('Ctrl+Shift+Space');
}

/** Rodziny, których drogi pilnuje ten plik, porównane z wykazem komend, które źródła naprawdę wysłały do rdzenia. */
const RODZINY = [
  'speech.',
  'memory.context.',
  'memory.retention.',
  'context.usage.',
  'clipboard.',
  'snippet.',
  'launcher.',
  'assistant.',
] as const;

describe('źródła rodzin przekrojowych modułu Assistant', () => {
  it('mają drogę z okna do każdej komendy swoich rodzin', async () => {
    const { kanal, wyslane } = kanalProbny();
    await przelotZrodel(kanal);

    const zRodzin = Object.values(Command).filter((nazwa) =>
      RODZINY.some((przedrostek) => nazwa.startsWith(przedrostek)),
    );
    const bezDrogi = zRodzin.filter((nazwa) => !wyslane.includes(nazwa));

    expect(bezDrogi, `komendy bez drogi z okna: ${bezDrogi.join(', ')}`).toEqual([]);
  });

  it('nie wysyła pól opcjonalnych, których Operator nie wskazał', async () => {
    const wyslane: Record<string, unknown> = {};
    const kanal = {
      wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
        wyslane[komenda] = zadanie;
        przyWyniku?.({ udany: true, wynik: {} });
        return 'zadanie-1';
      },
      naZdarzenie: () => () => undefined,
      naDowolny: () => () => undefined,
      sesja: () => ({}) as ReturnType<Kanal['sesja']>,
      dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
    } as unknown as Kanal;

    const schowek = utworzZrodloSchowka(kanal);
    await schowek.dopisz({ tresc: 'ustalenia' });
    const zadanie = wyslane[Command.ClipboardPush] as Record<string, unknown>;

    // Rodzaj wpisu, okno źródłowe i znacznik wrażliwości są rozstrzygnięciami Operatora.
    expect(Object.keys(zadanie).sort()).toEqual(['content']);

    const mowa = utworzZrodloMowy(kanal);
    await mowa.zapiszWybudzanie({ fraza: 'Danaco' });
    const nastawa = wyslane[Command.SpeechWakeSet] as Record<string, unknown>;

    // Zapis wybudzenia zostawia pola pominięte bez zmian, inaczej nadpisałby nastawę Operatora.
    expect(Object.keys(nastawa).sort()).toEqual(['phrase']);
  });
});
