import type { BrowserSnapshot } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { MonitorZmian } from './material-sesji';
import { skutekPobrania, skutekPrzechwycenia, skutekPrzekazania } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Rozmowa Capture & Monitor Panel z rdzeniem: przechwycenie strony, założenie
 * i sprawdzenie monitora zmian, otwarcie pozycji oraz przekazanie materiału
 * do modułów Library i Research.
 *
 * Jedna odpowiedzialność: wywołania rdzenia w imieniu panelu. Panel składa
 * kontrolki i wykaz, pamięć materiału stoi w `material-sesji.ts`, ocena
 * odpowiedzi w `skutek-zapisu.ts`.
 *
 * Sprawdzenie monitora przechodzi pod jego adres, bo `browser.snapshot.get`
 * oddaje treść strony bieżącej okna, a nie dowolnej. Znaczy to, że sprawdzenie
 * przestawia wspólny podgląd — Operator ma o tym wiedzieć przed naciśnięciem,
 * nie po nim, więc mówi mu to i dymek przy pozycji, i zdanie odpowiedzi.
 */
export interface CzynnosciMaterialu {
  /** `browser.snapshot.get` ze zrzutem i źródłem strony; pozycja materiału sesji. */
  przechwyc(): Promise<void>;
  /** Zakłada monitor na stronie widocznej w Browser Window. */
  zalozMonitor(): void;
  /** `browser.navigate` pod adres monitora i zestawienie treści z odniesieniem. */
  sprawdzMonitor(monitor: MonitorZmian): Promise<void>;
  /** `browser.navigate` pod adres pozycji materiału. */
  otworz(migawka: BrowserSnapshot): Promise<void>;
  /** `context.transfer` wykazu materiału do wskazanego modułu. */
  przekaz(kodModulu: string): Promise<void>;
}

export function utworzCzynnosciMaterialu(
  stan: StanPrzegladania,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): CzynnosciMaterialu {
  return {
    async przechwyc() {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return;
      }
      powiedz('Przechwycenie strony w toku…', true);
      const wynik = await stan.zrodlo.migawka({
        windowId: idOkna,
        includeScreenshot: true,
        includeHtml: true,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Przechwycenie strony', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      const migawka = wynik.wynik.snapshot;
      stan.wchlonMigawke(migawka);
      const skutek = skutekPrzechwycenia(migawka, stan.material.dopisz(migawka));
      powiedz(skutek.zdanie, skutek.udany);
    },

    zalozMonitor() {
      const migawka = stan.migawka();
      if (migawka === null) {
        powiedz('Monitor zakłada się na stronie widocznej — najpierw przejdź pod adres.', false);
        return;
      }
      if (!stan.material.zalozMonitor(migawka)) {
        powiedz(`Monitor strony ${migawka.url} już stoi — drugi pilnowałby tego samego.`, false);
        return;
      }
      powiedz(
        `Monitor założony na ${migawka.url}; odniesieniem jest treść z tej chwili ` +
          `(${(migawka.text ?? '').length} znaków).`,
        true,
      );
    },

    async sprawdzMonitor(monitor) {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return;
      }
      powiedz(`Sprawdzenie monitora ${monitor.adres} — wspólny podgląd przechodzi pod ten adres…`, true);
      const wynik = await stan.zrodlo.przejdz({ windowId: idOkna, url: monitor.adres });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Sprawdzenie monitora', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      const migawka = wynik.wynik.snapshot;
      stan.wchlonMigawke(migawka);
      const pobranie = skutekPobrania(migawka, monitor.adres, 'stronę monitora');
      if (!pobranie.udany) {
        powiedz(pobranie.zdanie, false);
        return;
      }
      const po = stan.material.zapiszSprawdzenie(monitor.adres, migawka);
      if (po === null) {
        powiedz('Monitor zniknął z wykazu w trakcie sprawdzenia — wynik nie ma gdzie stanąć.', false);
        return;
      }
      powiedz(zdanieOSprawdzeniu(po), true);
    },

    async otworz(migawka) {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return;
      }
      powiedz(`Otwieranie strony ${migawka.url}…`, true);
      const wynik = await stan.zrodlo.przejdz({ windowId: idOkna, url: migawka.url });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Otwarcie pozycji', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      stan.wchlonMigawke(wynik.wynik.snapshot);
      const skutek = skutekPobrania(wynik.wynik.snapshot, migawka.url, 'stronę');
      powiedz(skutek.zdanie, skutek.udany);
    },

    async przekaz(kodModulu) {
      const idOkna = stan.idOkna();
      const pozycje = stan.material.przechwycenia();
      if (idOkna === '' || pozycje.length === 0) {
        powiedz('Nie ma czego przekazać — najpierw przechwyć co najmniej jedną stronę.', false);
        return;
      }
      powiedz(`Przekazanie materiału sesji do modułu ${kodModulu}…`, true);
      // Komplet kontekstu nie ma pola na migawki strony, więc wykaz idzie
      // poleceniem wyjściowym. Zdanie o skutku mówi to wprost — inaczej
      // Operator sądziłby, że w module docelowym stanął obraz, a nie spis.
      const wykaz = pozycje
        .map((pozycja) => `${pozycja.migawka.url} (migawka ${pozycja.migawka.id})`)
        .join('\n');
      const wynik = await stan.zapisy.przekaz(idOkna, kodModulu, {
        prompt: `Materiał sesji przeglądania — ${pozycje.length} pozycji:\n${wykaz}`,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(
          opisOdmowy(`Przekazanie do ${kodModulu}`, wynik.blad?.code, wynik.blad?.message),
          false,
        );
        return;
      }
      const skutek = skutekPrzekazania(
        wynik.wynik,
        kodModulu,
        `Wysłano spis ${pozycje.length} pozycji materiału (adresy i identyfikatory migawek, bez obrazów)`,
      );
      powiedz(skutek.zdanie, skutek.udany);
    },
  };
}

/** Zdanie o wyniku sprawdzenia monitora wraz z miarą różnicy. */
function zdanieOSprawdzeniu(monitor: MonitorZmian): string {
  if (monitor.wynik === 'bez-zmian') {
    return `Monitor ${monitor.adres}: treść bez zmian względem odniesienia.`;
  }
  const kierunek = monitor.roznicaZnakow > 0 ? 'dłuższa' : 'krótsza';
  return (
    `Monitor ${monitor.adres}: treść ZMIENIONA — ${kierunek} o ` +
    `${Math.abs(monitor.roznicaZnakow)} znaków względem odniesienia. ` +
    'Odniesienie zostaje sprzed zmiany, więc kolejne sprawdzenie porówna się z nim samym.'
  );
}
