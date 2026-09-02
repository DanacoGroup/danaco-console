// Panele uniwersalne okna Library: miary repozytorium, artefakty, dziennik,
// konsola higieny i udostępnienia. Prototyp nie ma tu przycisków czynności,
// więc operacje niesie wskazanie wiersza, a naddatek — klawisz Shift.
import {
  Command,
  LibraryPackageKind,
  LibraryPreservationKind,
  LibraryShareScope,
  type LibraryRetentionPolicy,
  type LibraryShare,
  type LibraryStats,
  type LibraryWebhook,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  dzien,
  etykietaGrupy,
  panel,
  pozycja,
  pustka,
  rozmiar,
  wartosc,
  zdejmijDzieci,
  type Kontekst,
} from './library-wspolne.ts';

const GRANICA_DZIENNIKA = 50;

export function zwiazZadania(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  zwiazMiary(kanal, korzen, kontekst, przy);
  zwiazArtefakty(kanal, korzen, kontekst, przy);
  zwiazDziennik(kanal, korzen, kontekst);
  zwiazKonsole(kanal, korzen, kontekst);
  zwiazUdostepnienia(kanal, korzen, kontekst, przy);
}

function zwiazMiary(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const trescMozeMoze = panel(korzen, 'panel-zadania');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const wzorPostepu = tresc.querySelector('.dn-postep');
  const postep = wzorPostepu instanceof HTMLElement ? wzorPostepu : null;
  const polityki = new Map<HTMLElement, LibraryRetentionPolicy>();

  async function odswiez(): Promise<void> {
    const [odpowiedzMiar, odpowiedzRetencji] = await Promise.all([
      wywolaj(kanal, Command.LibraryStatsGet, {}),
      wywolaj(kanal, Command.LibraryRetentionList, {}),
    ]);
    const miary = wartosc(odpowiedzMiar)?.stats ?? null;
    const wykaz = wartosc(odpowiedzRetencji)?.policies ?? [];
    polityki.clear();
    if (miary === null && wykaz.length === 0) {
      pustka(tresc, 'Rdzeń nie podał miar repozytorium.');
      return;
    }
    zdejmijDzieci(tresc);
    if (miary !== null) naniesMiary(tresc, miary, postep);
    if (wykaz.length === 0) return;
    tresc.appendChild(etykietaGrupy(tresc.ownerDocument, 'Retencja'));
    for (const polityka of wykaz) {
      const wiersz = pozycja(
        tresc.ownerDocument,
        `${polityka.action} po ${String(polityka.keepDays)} dniach`,
        polityka.scopeId ?? polityka.scope,
      );
      polityki.set(wiersz, polityka);
      tresc.appendChild(wiersz);
    }
  }

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      if (!zdarzenie.shiftKey) return;
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wiersz = cel.closest('.dn-wykaz-modulu-poz');
      const polityka = wiersz instanceof HTMLElement ? polityki.get(wiersz) : undefined;
      if (polityka === undefined) return;
      void zdejmijRetencje(kanal, polityka).then(odswiez);
    },
    przy,
  );

  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function naniesMiary(tresc: HTMLElement, miary: LibraryStats, postep: HTMLElement | null): void {
  const dokument = tresc.ownerDocument;
  tresc.appendChild(pozycja(dokument, 'Pliki w repozytorium', String(miary.fileCount)));
  tresc.appendChild(pozycja(dokument, 'W archiwum', String(miary.archivedCount)));
  tresc.appendChild(pozycja(dokument, 'Bez przypisania', String(miary.orphanCount)));
  tresc.appendChild(pozycja(dokument, 'Duplikaty', String(miary.duplicateCount)));
  tresc.appendChild(pozycja(dokument, 'Zajętość', rozmiar(miary.totalBytes)));
  if (postep === null) return;
  const zsumowane = miary.fileCount - miary.missingChecksumCount;
  const udzial = miary.fileCount === 0 ? 0 : Math.round((zsumowane / miary.fileCount) * 100);
  postep.dataset.postepDo = String(udzial);
  const wartoscPaska = postep.querySelector<HTMLElement>('.dn-postep-wartosc');
  if (wartoscPaska !== null) wartoscPaska.style.width = `${String(udzial)}%`;
  tresc.appendChild(postep);
}

function zwiazArtefakty(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const trescMozeMoze = panel(korzen, 'panel-artefakty');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;

  async function odswiez(): Promise<void> {
    const odpowiedz = await wywolaj(kanal, Command.LibraryThesaurusExport, {
      includeCollections: true,
    });
    const slownik = wartosc(odpowiedz);
    if (slownik === null) {
      pustka(tresc, 'Rdzeń nie zwrócił artefaktów repozytorium.');
      return;
    }
    zdejmijDzieci(tresc);
    const opis = `${String(slownik.conceptCount)} pojęć · ${String(slownik.relationCount)} relacji`;
    const wiersz = pozycja(tresc.ownerDocument, `slownik-etykiet.${slownik.format}`, opis);
    wiersz.dataset.artefakt = 'slownik';
    tresc.appendChild(wiersz);
  }

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element) || cel.closest('.dn-wykaz-modulu-poz') === null) return;
      const plik = kontekst.plik();
      if (zdarzenie.shiftKey) {
        void utrwal(kanal, plik === null ? [] : [plik.id]);
        return;
      }
      void spakuj(kanal, plik === null ? [] : [plik.id]);
    },
    przy,
  );

  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function zwiazDziennik(kanal: Kanal, korzen: ParentNode, kontekst: Kontekst): void {
  const trescMozeMoze = panel(korzen, 'panel-subagenci');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;

  async function odswiez(): Promise<void> {
    const plik = kontekst.plik();
    const odpowiedz = await wywolaj(kanal, Command.LibraryAuditList, {
      ...(plik === null ? {} : { fileId: plik.id }),
      limit: GRANICA_DZIENNIKA,
    });
    const wpisy = wartosc(odpowiedz)?.entries ?? [];
    if (wpisy.length === 0) {
      pustka(tresc, 'Rdzeń nie zwrócił czynności dla tego zakresu.');
      return;
    }
    zdejmijDzieci(tresc);
    for (const wpis of wpisy) {
      tresc.appendChild(
        pozycja(
          tresc.ownerDocument,
          `${wpis.action} — ${wpis.detail ?? wpis.actor}`,
          dzien(wpis.at),
        ),
      );
    }
  }

  kontekst.przyPliku(() => {
    void odswiez();
  });
  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function zwiazKonsole(kanal: Kanal, korzen: ParentNode, kontekst: Kontekst): void {
  const trescMoze = korzen.querySelector('#panel-terminal .pt-konsola');
  if (!(trescMoze instanceof HTMLElement)) return;
  const tresc = trescMoze;

  async function odswiez(): Promise<void> {
    const [odpowiedzSum, odpowiedzNazw] = await Promise.all([
      wywolaj(kanal, Command.LibraryFixityCheck, {}),
      wywolaj(kanal, Command.LibraryNameNormalize, { dryRun: true, transliterate: true }),
    ]);
    const sumy = wartosc(odpowiedzSum);
    const nazwy = wartosc(odpowiedzNazw);
    const wiersze: string[] = [];
    if (sumy !== null) {
      wiersze.push(
        `sumy kontrolne: sprawdzono ${String(sumy.checkedCount)}, rozbieżnych ${String(sumy.mismatchedCount)}, brakujących ${String(sumy.missingCount)}`,
      );
    }
    if (nazwy !== null) {
      wiersze.push(`normalizacja nazw (próba): do zmiany ${String(nazwy.changedCount)}`);
    }
    if (wiersze.length === 0) {
      tresc.textContent = 'Rdzeń nie zwrócił wyniku higieny repozytorium.';
      return;
    }
    tresc.textContent = wiersze.join('\n');
  }

  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function zwiazUdostepnienia(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const trescMozeMoze = panel(korzen, 'panel-przegladarka');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const pole = tresc.querySelector('input[type="search"]');
  const stopka = tresc.querySelector('.lb-mono');
  const udostepnienia = new Map<HTMLElement, LibraryShare>();
  let nasluchy: LibraryWebhook[] = [];

  async function odswiez(): Promise<void> {
    const [odpowiedzUdostepnien, odpowiedzNasluchow] = await Promise.all([
      wywolaj(kanal, Command.LibraryShareList, {}),
      wywolaj(kanal, Command.LibraryWebhookList, {}),
    ]);
    const wykaz = wartosc(odpowiedzUdostepnien)?.shares ?? [];
    nasluchy = wartosc(odpowiedzNasluchow)?.webhooks ?? [];
    udostepnienia.clear();
    for (const wiersz of tresc.querySelectorAll('.dn-wykaz-modulu-poz')) wiersz.remove();
    for (const udostepnienie of wykaz) {
      const wiersz = pozycja(
        tresc.ownerDocument,
        `${udostepnienie.scope} · ${udostepnienie.targetId}`,
        opiszWaznosc(udostepnienie),
      );
      udostepnienia.set(wiersz, udostepnienie);
      if (stopka === null) tresc.appendChild(wiersz);
      else tresc.insertBefore(wiersz, stopka);
    }
    if (stopka !== null) stopka.textContent = opiszNasluchy(nasluchy);
  }

  pole?.addEventListener(
    'keydown',
    (zdarzenie) => {
      if (!(zdarzenie instanceof KeyboardEvent) || zdarzenie.key !== 'Enter') return;
      const plikMozeMoze = kontekst.plik();
      if (plikMozeMoze === null) return;
      const plikMoze = plikMozeMoze;
      const plik = plikMoze;
      void wystaw(kanal, plik.id, pole).then(odswiez);
    },
    przy,
  );

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      if (cel.closest('.lb-mono') !== null) {
        void przestawNasluch(kanal, nasluchy, zdarzenie.shiftKey).then(odswiez);
        return;
      }
      const wiersz = cel.closest('.dn-wykaz-modulu-poz');
      const udostepnienie = wiersz instanceof HTMLElement ? udostepnienia.get(wiersz) : undefined;
      if (udostepnienie === undefined || !zdarzenie.shiftKey) return;
      void odwolaj(kanal, udostepnienie.id).then(odswiez);
    },
    przy,
  );

  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function opiszWaznosc(udostepnienie: LibraryShare): string {
  if (udostepnienie.revokedAt !== undefined) return 'odwołane';
  if (udostepnienie.expiresAt === undefined) return 'bez terminu';
  return `wygasa ${dzien(udostepnienie.expiresAt)}`;
}

function opiszNasluchy(nasluchy: LibraryWebhook[]): string {
  if (nasluchy.length === 0) return 'Rdzeń nie zwrócił nasłuchów.';
  const czynne = nasluchy.filter((nasluch) => nasluch.enabled).length;
  return `Nasłuchy: ${String(czynne)} czynnych z ${String(nasluchy.length)}`;
}

async function zdejmijRetencje(kanal: Kanal, polityka: LibraryRetentionPolicy): Promise<void> {
  await wywolaj(kanal, Command.LibraryRetentionSet, { policy: polityka, remove: true });
}

async function utrwal(kanal: Kanal, idPlikow: string[]): Promise<void> {
  if (idPlikow.length === 0) return;
  const odpowiedz = await wywolaj(kanal, Command.LibraryPreservationRun, {
    fileIds: idPlikow,
    kind: LibraryPreservationKind.Pdfa,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Utrwalono postaci: ${String(wynik.validCount)}.`);
}

async function spakuj(kanal: Kanal, idPlikow: string[]): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryPackageExport, {
    kind: LibraryPackageKind.Snapshot,
    ...(idPlikow.length === 0 ? {} : { fileIds: idPlikow }),
    includeVersions: true,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Pakiet: ${String(wynik.entries)} pozycji, ${rozmiar(wynik.sizeBytes)}.`);
}

async function wystaw(kanal: Kanal, idPliku: string, pole: Element): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryShareCreate, {
    scope: LibraryShareScope.File,
    targetId: idPliku,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  if (pole instanceof HTMLInputElement) pole.value = wynik.url;
}

async function odwolaj(kanal: Kanal, idUdostepnienia: string): Promise<void> {
  await wywolaj(kanal, Command.LibraryShareRevoke, { shareId: idUdostepnienia });
}

async function przestawNasluch(
  kanal: Kanal,
  nasluchy: LibraryWebhook[],
  zdejmij: boolean,
): Promise<void> {
  const nasluch = nasluchy[0];
  if (nasluch === undefined) return;
  if (zdejmij) {
    await wywolaj(kanal, Command.LibraryWebhookRemove, { webhookId: nasluch.id });
    return;
  }
  await wywolaj(kanal, Command.LibraryWebhookSet, {
    webhook: { ...nasluch, enabled: !nasluch.enabled },
  });
}
