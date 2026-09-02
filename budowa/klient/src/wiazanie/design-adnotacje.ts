// Wersje kompozycji, adnotacje i obecność współpracujących — wszystkie trzy
// wykazy trafiają w panel roboczy, bo prototyp nie ma dla nich osobnych węzłów.

import { Command } from '../../../shared/contract.ts';
import { oglos } from './ogloszenie.ts';
import {
  chwila,
  poproszony,
  wykonany,
  type Kontekst,
  type Wiersz,
} from './design-wspolne.ts';
import { pokazWykaz } from './design-wykaz.ts';
import { stanPlanszy } from './design-plansza.ts';

export async function pokazWersje(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idPlanszy === '') {
    oglos('Design', 'Wersje dotyczą kompozycji, której okno jeszcze nie ma.', 'ostrzezenie');
    return;
  }
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignBoardVersionList, {
    boardId: kontekst.stan.idPlanszy,
    limit: 50,
  });
  const wersje = odpowiedz?.versions ?? [];
  const wiersze: Wiersz[] = wersje.map((wersja) => ({
    tekst: wersja.name ?? wersja.id,
    meta: `${String(wersja.layerCount)} warstw · ${chwila(wersja.createdAt)}`,
    plakietka: wersja.note ?? '',
    kropka: 'neutralna',
    naKlik: () => {
      void przywroc(kontekst, wersja.id);
    },
  }));
  wiersze.unshift({
    tekst: 'Zapisz obecny układ jako wersję',
    kropka: 'sygnal',
    naKlik: () => {
      void zapisz(kontekst);
    },
  });
  pokazWykaz(kontekst, 'Wersje kompozycji', wiersze);
}

async function zapisz(kontekst: Kontekst): Promise<void> {
  await wykonany(kontekst.kanal, Command.DesignBoardVersionSave, {
    boardId: kontekst.stan.idPlanszy,
    name: `wersja ${chwila(Date.now())}`,
  }, 'Zapis wersji');
  await pokazWersje(kontekst);
}

async function przywroc(kontekst: Kontekst, idWersji: string): Promise<void> {
  await wykonany(kontekst.kanal, Command.DesignBoardVersionRestore, {
    versionId: idWersji,
  }, 'Przywrócenie wersji');
  kontekst.stan.odswiezenia.get('plansza')?.();
}

export async function pokazAdnotacje(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idPlanszy === '') {
    oglos('Design', 'Adnotacja stoi na kompozycji, której okno jeszcze nie ma.', 'ostrzezenie');
    return;
  }
  await zglosObecnosc(kontekst);
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignAnnotationList, {
    boardId: kontekst.stan.idPlanszy,
  });
  const adnotacje = odpowiedz?.annotations ?? [];
  const wiersze: Wiersz[] = adnotacje.map((adnotacja) => ({
    tekst: adnotacja.text,
    meta: `${adnotacja.author ?? 'bez autora'} · ${chwila(adnotacja.createdAt)}`,
    kropka: adnotacja.resolved === true ? 'sukces' : 'sygnal',
    naKlik: () => {
      void domknij(kontekst, adnotacja.id, adnotacja.text, adnotacja.resolved !== true);
    },
  }));
  const stan = stanPlanszy(kontekst.idKarty);
  const zaznaczona = [...stan.zaznaczone][0];
  if (zaznaczona !== undefined) {
    wiersze.unshift({
      tekst: 'Dopisz adnotację do zaznaczonej warstwy',
      kropka: 'neutralna',
      naKlik: () => {
        void dopisz(kontekst, zaznaczona);
      },
    });
  }
  pokazWykaz(kontekst, 'Adnotacje kompozycji', wiersze, 'Kompozycja nie ma adnotacji.');
}

async function dopisz(kontekst: Kontekst, idWarstwy: string): Promise<void> {
  await wykonany(kontekst.kanal, Command.DesignAnnotationSet, {
    boardId: kontekst.stan.idPlanszy,
    layerId: idWarstwy,
    text: 'Do przejrzenia.',
  }, 'Adnotacja');
  await pokazAdnotacje(kontekst);
}

async function domknij(
  kontekst: Kontekst,
  idAdnotacji: string,
  tekst: string,
  domkniecie: boolean,
): Promise<void> {
  await wykonany(kontekst.kanal, Command.DesignAnnotationSet, {
    boardId: kontekst.stan.idPlanszy,
    annotationId: idAdnotacji,
    text: tekst,
    resolved: domkniecie,
  }, domkniecie ? 'Domknięcie adnotacji' : 'Otwarcie adnotacji');
  await pokazAdnotacje(kontekst);
}

// Zgłoszenie obecności jest ulotne: rdzeń rozgłasza je i nie zapisuje.
export async function zglosObecnosc(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idPlanszy === '') return;
  const stan = stanPlanszy(kontekst.idKarty);
  await poproszony(kontekst.kanal, Command.DesignPresenceReport, {
    boardId: kontekst.stan.idPlanszy,
    selectedLayerIds: [...stan.zaznaczone],
  });
}

export async function zejdzZKompozycji(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idPlanszy === '') return;
  await poproszony(kontekst.kanal, Command.DesignPresenceReport, {
    boardId: kontekst.stan.idPlanszy,
    leaving: true,
  });
}
