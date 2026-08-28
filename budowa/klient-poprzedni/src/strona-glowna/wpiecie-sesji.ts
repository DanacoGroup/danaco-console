import type { Kanal } from '../protokol/kanal';
import { zadajPrzeniesienieOgniska } from '../protokol/ognisko-sesji';
import { zadajPowiazanieSesji } from '../protokol/powiazanie-sesji';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import type { KodSrodowiska } from './pozycje-srodowisk';
import type { StronaGlowna } from './okno-strona-glowna';
import { utworzStrefeArchiwum } from './strefa-archiwum';
import { utworzCzynnosciHistorii, utworzOdpytanieArchiwum } from './wykonanie-czynnosci';
import { utworzZrodloSesji } from './zrodlo-sesji';

/** Wpięcie źródła sesji w stronę główną wiąże migawki historii ze stroną i wykonuje powrót do sesji trwającej na rdzeniu. */
export interface ZaleznosciWpieciaSesji {
  strona: StronaGlowna;
  /** Kanał kontraktu — `session.list`, zdarzenia i komendy powrotu. */
  kanal: Kanal;
  /** Tożsamość klienta z powitania; bez niej powrót pozostaje wyłączony. */
  klient?: TozsamoscKlienta;
  /** Przejście do środowiska sesji po udanym powrocie. */
  naPrzejscie?(kod: KodSrodowiska, modul?: string): void;
}

export function zasilSesjeStronyGlownej(zaleznosci: ZaleznosciWpieciaSesji): void {
  const { strona, kanal, klient } = zaleznosci;
  const zrodlo = utworzZrodloSesji(kanal);

  strona.sesje.ustawMigawke(zrodlo.migawka());
  zrodlo.naZmiane((migawka) => strona.sesje.ustawMigawke(migawka));

  // Archiwum odpytuje rdzeń samo przy rozwinięciu; odświeżenie domyka pętlę po innych czynnościach.
  const czynnosci = utworzCzynnosciHistorii({ kanal, odswiez: () => archiwum.odswiez() });
  const archiwum = utworzStrefeArchiwum(
    utworzOdpytanieArchiwum(kanal),
    czynnosci,
    (tekst) => strona.sesje.zglosKomunikat(tekst),
    strona.wykazSrodowisk,
  );
  strona.sesje.ustawCzynnosci(czynnosci);
  strona.sesje.osadzArchiwum(archiwum.element);

  if (klient === undefined) return;

  strona.sesje.ustawPowrot(async (wpis) => {
    const powiazanie = await zadajPowiazanieSesji(kanal, {
      sessionId: wpis.sesja.id,
      clientId: klient.id,
    });
    if (!powiazanie.udany || powiazanie.wynik?.bound !== true) {
      strona.sesje.zglosKomunikat(
        powiazanie.blad?.message ?? 'Rdzeń nie powiązał połączenia z sesją.',
      );
      return;
    }

    // Od tej chwili koperty wychodzące niosą sesję powiązaną, tak jak po nawiązaniu połączenia.
    kanal.sesja().ustaw(wpis.sesja.id);

    const ognisko = await zadajPrzeniesienieOgniska(kanal, {
      sessionId: wpis.sesja.id,
      clientId: klient.id,
      windowId: wpis.obecnosc?.focusedWindowId,
    });
    if (!ognisko.udany) {
      // Powiązanie już trwa; odmowa ogniska nie przerywa powrotu.
      console.warn('[strona główna] ognisko po powrocie nie przeszło', ognisko.blad?.message);
    }

    // Kod rozstrzyga wykaz z rdzenia, nie stałą kliencką, by powrót prowadził od razu do pracy.
    const kod = strona.wykazSrodowisk.kod(wpis.obecnosc?.environmentCode);
    if (kod !== null) zaleznosci.naPrzejscie?.(kod, wpis.obecnosc?.moduleCode);
  });
}
