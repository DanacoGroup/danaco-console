// Szyna modułu Terminal: hosty zdalne, klucze, tunele i obserwacje plików.
// Pozycje szyny stawiane są z wykazów rdzenia, nie z wzoru w znaczniku.

import {
  Command,
  type TerminalHost,
  type TerminalSshKey,
  type TerminalTunnel,
  type TerminalWatch,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

export async function zwiazWyposazenieTerminala(
  kanal: Kanal,
  korzen: Element,
  przy: AddEventListenerOptions,
): Promise<void> {
  const szyna = korzen.querySelector('.dn-szyna-modulu-lista');
  if (szyna === null) return;
  zdejmijPozycjeWzorcowe(szyna);

  await Promise.all([
    postawHosty(kanal, szyna),
    postawKlucze(kanal, szyna),
    postawTunele(kanal, szyna),
    postawObserwacje(kanal, szyna),
  ]);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const tunel = cel.closest<HTMLElement>('[data-tunel-zamknij]');
    if (tunel !== null) {
      void zamknijTunel(kanal, szyna, tunel.dataset.tunelZamknij ?? '');
      return;
    }
    const obserwacja = cel.closest<HTMLElement>('[data-obserwacja-zdejmij]');
    if (obserwacja !== null) {
      void zdejmijObserwacje(kanal, szyna, obserwacja.dataset.obserwacjaZdejmij ?? '');
    }
  }, przy);
}

/* Znacznik niesie nazwy hostów i zestawów wpisane wprost; wypełniacz zostaje
   w produkcie na zawsze, więc pozycje wzorcowe schodzą przed wypełnieniem. */
function zdejmijPozycjeWzorcowe(szyna: Element): void {
  for (const pozycja of szyna.querySelectorAll('.pt-pozycja')) pozycja.remove();
}

async function postawHosty(kanal: Kanal, szyna: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalHostList, {});
  if (!wynik.udany) return;
  const hosty = (wynik.wynik as { hosts?: TerminalHost[] } | undefined)?.hosts ?? [];
  postawGrupe(szyna, 'Hosty zdalne', hosty.map((host) => ({
    tytul: host.name,
    cechy: { hostTerminala: host.id },
  })));
}

async function postawKlucze(kanal: Kanal, szyna: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalKeyList, {});
  if (!wynik.udany) return;
  const klucze = (wynik.wynik as { keys?: TerminalSshKey[] } | undefined)?.keys ?? [];
  postawGrupe(szyna, 'Klucze', klucze.map((klucz) => ({
    tytul: klucz.name,
    cechy: { kluczTerminala: klucz.id },
  })));
}

async function postawTunele(kanal: Kanal, szyna: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalTunnelList, {});
  if (!wynik.udany) return;
  const tunele = (wynik.wynik as { tunnels?: TerminalTunnel[] } | undefined)?.tunnels ?? [];
  postawGrupe(szyna, 'Tunele', tunele.map((tunel) => ({
    tytul: `${tunel.hostId ?? tunel.kind} → ${String(tunel.remotePort ?? 0)}`,
    cechy: { tunelZamknij: tunel.id },
  })));
}

async function postawObserwacje(kanal: Kanal, szyna: Element): Promise<void> {
  const wynik = await wywolaj(kanal, Command.TerminalWatchList, {});
  if (!wynik.udany) return;
  const obserwacje = (wynik.wynik as { watches?: TerminalWatch[] } | undefined)?.watches ?? [];
  postawGrupe(szyna, 'Obserwacje', obserwacje.map((obserwacja) => ({
    tytul: obserwacja.pattern,
    cechy: { obserwacjaZdejmij: obserwacja.id },
  })));
}

interface PozycjaSzyny {
  tytul: string;
  cechy: Record<string, string>;
}

function postawGrupe(szyna: Element, etykieta: string, pozycje: PozycjaSzyny[]): void {
  if (pozycje.length === 0) return;
  const naglowek = document.createElement('div');
  naglowek.className = 'dn-etyk-mono';
  naglowek.textContent = etykieta;
  szyna.append(naglowek);
  for (const pozycja of pozycje) {
    const przycisk = document.createElement('button');
    przycisk.className = 'pt-pozycja';
    przycisk.type = 'button';
    for (const [nazwa, wartosc] of Object.entries(pozycja.cechy)) {
      przycisk.dataset[nazwa] = wartosc;
    }
    const tytul = document.createElement('span');
    tytul.className = 'pt-pozycja-tytul';
    tytul.textContent = pozycja.tytul;
    przycisk.append(tytul);
    szyna.append(przycisk);
  }
}

async function zamknijTunel(kanal: Kanal, szyna: Element, id: string): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.TerminalTunnelClose, { tunnelId: id });
  if (!wynik.udany) {
    oglos('Terminal', wynik.blad?.message ?? 'Tunel nie został zamknięty.', 'blad');
    return;
  }
  void postawTunele(kanal, szyna);
}

async function zdejmijObserwacje(kanal: Kanal, szyna: Element, id: string): Promise<void> {
  if (id === '') return;
  const wynik = await wywolaj(kanal, Command.TerminalWatchStop, { watchId: id });
  if (!wynik.udany) {
    oglos('Terminal', wynik.blad?.message ?? 'Obserwacja nie została zdjęta.', 'blad');
    return;
  }
  void postawObserwacje(kanal, szyna);
}
