/**
 * Panel Session Repository okna Studia — repozytorium sesji dokumentu w czterech
 * widokach: wersje, gałęzie, dziennik czynności i kopie zapasowe. Znacznik i
 * klasy wzięte ze źródła kształtu (`design/05-okna/moduly/studio.html`,
 * `section#panel-repo`); wiersz wersji stoi na klasach `dn-wersja-*`.
 *
 * Rodziny komend tego panelu (`studio.repository.*`, `studio.version.*`,
 * `studio.branch.*`, `studio.journal.*`, `studio.backup.*`) operują na
 * `documentId`, którego `ZaleznosciPanelu` nie niesie — panel poznaje go
 * wyłącznie ze zdarzenia `studio.document.changed` ograniczonego do własnego
 * okna. Do pierwszego takiego zdarzenia panel czeka, nie zgaduje dokumentu.
 *
 * Podgląd i porównanie wersji zostają poza tym plikiem: prowadzą je okna
 * Preview Window i Diff/Grep Panel. Strefa nazywa to zdaniem zamiast stawiać
 * przyciski, które nie mają dokąd prowadzić.
 */

import {
  ChangeKind,
  Command,
  EventType,
  StudioActionState,
  StudioAuthor,
  StudioBackupReason,
  StudioMergeSide,
  StudioVersionSeries,
  type ErrorInfo,
  type StudioActionBalance,
  type StudioActionKind,
  type StudioBranch,
  type StudioDocumentAction,
  type StudioDocumentBackup,
  type StudioMergeConflict,
  type StudioMergeResolution,
  type StudioVersion,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import type { MontazPanelu, ZamontowanyPanel, ZaleznosciPanelu } from './umowa.ts';
import {
  komunikatCofniecia,
  komunikatEksportu,
  komunikatGalezi,
  komunikatKonfliktow,
  komunikatKopiiNowyDokument,
  komunikatOdwolania,
  komunikatPonowienia,
  komunikatScalenia,
  opisKonfliktu,
  podpisPunktuStartowego,
  podpisWersji,
  opisNiezapisanych,
  opisSzeregow,
  opisWpisow,
  tresciRepo,
  znacznikWersji,
} from './repo-tresci.ts';
import { opisOdmowy, zaloguj } from '../odmowa.ts';

type Karta = 'wersje' | 'galezie' | 'dziennik' | 'kopie';

type StanPodstawy = { rodzaj: 'brakOkna' } | { rodzaj: 'oczekiwanieDokumentu' } | { rodzaj: 'dokument'; id: string };

type StanWersji =
  | null
  | { rodzaj: 'wczytywanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'gotowe'; wersje: StudioVersion[]; operatorow?: number; autozapisow?: number; zalozycielska?: string };

type StanGalezi =
  | null
  | { rodzaj: 'wczytywanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'gotowe'; galezie: StudioBranch[] };

type StanDziennika =
  | null
  | { rodzaj: 'wczytywanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'gotowe'; wpisy: StudioDocumentAction[]; wszystkich: number };

type StanKopii =
  | null
  | { rodzaj: 'wczytywanie' }
  | { rodzaj: 'odmowa'; blad?: ErrorInfo }
  | { rodzaj: 'gotowe'; kopie: StudioDocumentBackup[]; niezapisanych: number };

type RodzajFormularza = 'rozgalez' | 'etykieta';
type StanFormularza = { akcja: RodzajFormularza; wersjaId: string } | null;
type StanKomunikatu = { rodzaj: 'sukces' | 'blad'; tekst: string } | null;
type StanAkcjiWiersza = { wersjaId: string; akcja: 'przywroc' | 'odwolanie' | RodzajFormularza } | null;
type StanScalania = { zrodlo: string; cel: string; konflikty: StudioMergeConflict[]; wybory: Map<number, StudioMergeSide> } | null;

type WyborSzeregu = 'wszystkie' | StudioVersionSeries;
type WyborAutora = 'wszyscy' | StudioAuthor;
type WyborStanu = 'wszystkie' | StudioActionState;

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function opisAutora(autor: StudioAuthor | undefined): string {
  if (autor === StudioAuthor.Model) return tresciRepo.wiersz.autorModel;
  if (autor === StudioAuthor.Uzytkownik) return tresciRepo.wiersz.autorUzytkownik;
  return tresciRepo.wiersz.autorNieznany;
}

function opisRodzaju(rodzaj: StudioActionKind): string {
  return tresciRepo.rodzajCzynnosci[rodzaj] ?? tresciRepo.rodzajCzynnosci.nieznany;
}

function opisPowodu(powod: StudioBackupReason): string {
  switch (powod) {
    case StudioBackupReason.Interval:
      return tresciRepo.kopie.powodInterval;
    case StudioBackupReason.Event:
      return tresciRepo.kopie.powodEvent;
    case StudioBackupReason.BeforeIrreversible:
      return tresciRepo.kopie.powodPrzedNieodwracalna;
    case StudioBackupReason.Manual:
      return tresciRepo.kopie.powodReczny;
    default:
      return tresciRepo.rodzajCzynnosci.nieznany;
  }
}

function formatCzas(znacznikCzasu: number): string {
  return new Date(znacznikCzasu).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

function formatBajty(bajty: number): string {
  if (bajty < 1024) return `${bajty} B`;
  const kb = bajty / 1024;
  if (kb < 1024) return `${kb.toFixed(1)} kB`;
  return `${(kb / 1024).toFixed(1)} MB`;
}

export const montujPanelRepo: MontazPanelu = (wezel, zaleznosci) => zamontuj(wezel, zaleznosci);

function zamontuj(wezel: HTMLElement, zaleznosci: ZaleznosciPanelu): ZamontowanyPanel {
  wezel.classList.add('sta-okno');
  wezel.id = 'panel-repo';
  wezel.setAttribute('role', 'tabpanel');
  wezel.setAttribute('aria-labelledby', 'karta-repo');
  wezel.setAttribute('data-nazwa', tresciRepo.panel.tytul);

  let zdjete = false;
  let podstawa: StanPodstawy = zaleznosci.idOkna === null ? { rodzaj: 'brakOkna' } : { rodzaj: 'oczekiwanieDokumentu' };
  let karta: Karta = 'wersje';

  let stanWersji: StanWersji = null;
  let stanGalezi: StanGalezi = null;
  let stanDziennika: StanDziennika = null;
  let stanKopii: StanKopii = null;

  let szereg: WyborSzeregu = 'wszystkie';
  let tylkoKluczowe = false;
  let autorWpisow: WyborAutora = 'wszyscy';
  let stanWpisow: WyborStanu = 'wszystkie';
  let tylkoNiezapisane = false;

  let formularz: StanFormularza = null;
  let komunikat: StanKomunikatu = null;
  let akcjaWiersza: StanAkcjiWiersza = null;
  let eksportTrwa = false;
  let powrotTrwa = false;
  let scalanaGalaz: string | null = null;
  let scalenie: StanScalania = null;
  let czynnoscWBiegu: string | null = null;
  let zakladanieKopii = false;
  let kopiaWBiegu: string | null = null;

  const znacznik = el('span', { klasa: 'sta-okno-znacznik', hidden: true });
  const tresc = el('div', { klasa: 'sta-okno-tresc' });
  const naglowek = el('header', { klasa: 'sta-okno-belka' }, [
    el('span', { klasa: 'sta-okno-tytul' }, [znak('historia'), el('b', { tekst: tresciRepo.panel.tytul })]),
    znacznik,
  ]);
  wezel.replaceChildren(naglowek, tresc);

  /** Dokument okna albo brak — każda komenda tej strefy potrzebuje go w żądaniu. */
  function idDokumentu(): string | null {
    return podstawa.rodzaj === 'dokument' ? podstawa.id : null;
  }

  // ── Wczytywanie widoków ────────────────────────────────────────────────

  async function wczytajWersje(dokument: string): Promise<void> {
    stanWersji = { rodzaj: 'wczytywanie' };
    odswiez();
    /* Wykaz nierozdzielony bierze `repository.list`; zawężenie do szeregu albo
       do wersji kluczowych bierze `version.series.list`, bo tylko ono rozdziela
       zapisy samoczynne od wersji Operatora i podaje wersję założycielską. */
    if (szereg === 'wszystkie' && !tylkoKluczowe) {
      const wynik = await wywolaj(zaleznosci.kanal, Command.StudioRepositoryList, { documentId: dokument });
      if (zdjete) return;
      stanWersji =
        !wynik.udany || wynik.wynik === undefined
          ? { rodzaj: 'odmowa', blad: wynik.blad }
          : { rodzaj: 'gotowe', wersje: wynik.wynik.versions };
      odswiez();
      return;
    }
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioVersionSeriesList, {
      documentId: dokument,
      series: szereg === 'wszystkie' ? undefined : szereg,
      milestonesOnly: tylkoKluczowe ? true : undefined,
    });
    if (zdjete) return;
    stanWersji =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'odmowa', blad: wynik.blad }
        : {
            rodzaj: 'gotowe',
            wersje: wynik.wynik.versions,
            operatorow: wynik.wynik.operatorCount,
            autozapisow: wynik.wynik.autosaveCount,
            zalozycielska: wynik.wynik.initialVersionId,
          };
    odswiez();
  }

  async function wczytajGalezie(dokument: string): Promise<void> {
    stanGalezi = { rodzaj: 'wczytywanie' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBranchList, { documentId: dokument });
    if (zdjete) return;
    stanGalezi =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'odmowa', blad: wynik.blad }
        : { rodzaj: 'gotowe', galezie: wynik.wynik.branches };
    odswiez();
  }

  async function wczytajDziennik(dokument: string): Promise<void> {
    stanDziennika = { rodzaj: 'wczytywanie' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioJournalList, {
      documentId: dokument,
      author: autorWpisow === 'wszyscy' ? undefined : autorWpisow,
      state: stanWpisow === 'wszystkie' ? undefined : stanWpisow,
    });
    if (zdjete) return;
    stanDziennika =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'odmowa', blad: wynik.blad }
        : { rodzaj: 'gotowe', wpisy: wynik.wynik.actions, wszystkich: wynik.wynik.total };
    odswiez();
  }

  async function wczytajKopie(dokument: string): Promise<void> {
    stanKopii = { rodzaj: 'wczytywanie' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBackupList, {
      documentId: dokument,
      windowId: zaleznosci.idOkna ?? undefined,
      unsavedOnly: tylkoNiezapisane ? true : undefined,
    });
    if (zdjete) return;
    stanKopii =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'odmowa', blad: wynik.blad }
        : { rodzaj: 'gotowe', kopie: wynik.wynik.backups, niezapisanych: wynik.wynik.unsavedCount };
    odswiez();
  }

  /** Wczytuje widok, na który Operator właśnie patrzy — reszta czeka do przełączenia. */
  function wczytajBiezacy(przeladuj: boolean): void {
    const dokument = idDokumentu();
    if (dokument === null) return;
    if (karta === 'wersje' && (przeladuj || stanWersji === null)) void wczytajWersje(dokument);
    if (karta === 'galezie' && (przeladuj || stanGalezi === null)) void wczytajGalezie(dokument);
    if (karta === 'dziennik' && (przeladuj || stanDziennika === null)) void wczytajDziennik(dokument);
    if (karta === 'kopie' && (przeladuj || stanKopii === null)) void wczytajKopie(dokument);
  }

  // ── Czynności na wersjach ──────────────────────────────────────────────

  async function przywroc(dokument: string, wersja: StudioVersion): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'przywroc' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioRepositoryRestore, {
      documentId: dokument,
      versionId: wersja.id,
    });
    if (zdjete) return;
    akcjaWiersza = null;
    if (!wynik.udany) {
      zaloguj(wynik.blad, 'repo.przywroc');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.przywroc };
      odswiez();
      return;
    }
    komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.przywrocono };
    await wczytajWersje(dokument);
  }

  async function zalozGalaz(dokument: string, wersja: StudioVersion, nazwa: string): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'rozgalez' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBranchCreate, {
      documentId: dokument,
      fromVersionId: wersja.id,
      name: nazwa,
    });
    if (zdjete) return;
    akcjaWiersza = null;
    formularz = null;
    zaloguj(wynik.blad, 'repo.galaz');
    if (!wynik.udany || wynik.wynik === undefined) {
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.galaz };
      odswiez();
      return;
    }
    komunikat = { rodzaj: 'sukces', tekst: komunikatGalezi(wynik.wynik.branch.name) };
    stanGalezi = null;
    odswiez();
  }

  async function ustawEtykiete(wersja: StudioVersion, etykieta: string, kluczowa: boolean): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'etykieta' };
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioVersionLabelSet, {
      versionId: wersja.id,
      label: etykieta === '' ? undefined : etykieta,
      milestone: kluczowa,
    });
    if (zdjete) return;
    akcjaWiersza = null;
    formularz = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.etykieta');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.etykieta };
      odswiez();
      return;
    }
    if (stanWersji !== null && stanWersji.rodzaj === 'gotowe') {
      const podmieniona = wynik.wynik.version;
      stanWersji = { ...stanWersji, wersje: stanWersji.wersje.map((w) => (w.id === podmieniona.id ? podmieniona : w)) };
    }
    komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.etykietaZapisana };
    odswiez();
  }

  async function utworzOdwolanie(wersja: StudioVersion): Promise<void> {
    akcjaWiersza = { wersjaId: wersja.id, akcja: 'odwolanie' };
    odswiez();
    /* Zasięg widoczności zostaje pusty: kontrakt bierze wtedy zasięg sesji,
       a strefa nie ma powierzchni, na której Operator wskazałby szerszy. */
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioVersionReferenceCreate, { versionId: wersja.id });
    if (zdjete) return;
    akcjaWiersza = null;
    zaloguj(wynik.blad, 'repo.odwolanie');
    komunikat =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'blad', tekst: tresciRepo.odmowa.odwolanie }
        : { rodzaj: 'sukces', tekst: komunikatOdwolania(wynik.wynik.reference) };
    odswiez();
  }

  async function eksportuj(dokument: string): Promise<void> {
    eksportTrwa = true;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioRepositoryExport, {
      documentId: dokument,
      windowId: zaleznosci.idOkna ?? undefined,
    });
    if (zdjete) return;
    eksportTrwa = false;
    zaloguj(wynik.blad, 'repo.eksport');
    komunikat =
      !wynik.udany || wynik.wynik === undefined
        ? { rodzaj: 'blad', tekst: tresciRepo.odmowa.eksport }
        : { rodzaj: 'sukces', tekst: komunikatEksportu(wynik.wynik.entries, formatBajty(wynik.wynik.sizeBytes)) };
    odswiez();
  }

  async function wrocDoZalozycielskiej(dokument: string): Promise<void> {
    powrotTrwa = true;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioVersionRestoreInitial, { documentId: dokument });
    if (zdjete) return;
    powrotTrwa = false;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.zalozycielska');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.zalozycielska };
      odswiez();
      return;
    }
    komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.wrocono };
    await wczytajWersje(dokument);
  }

  // ── Czynności na gałęziach ─────────────────────────────────────────────

  async function scal(zrodlo: StudioBranch, celId: string, rozstrzygniecia?: StudioMergeResolution[]): Promise<void> {
    scalanaGalaz = zrodlo.id;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBranchMerge, {
      sourceBranchId: zrodlo.id,
      targetBranchId: celId,
      resolutions: rozstrzygniecia,
    });
    if (zdjete) return;
    scalanaGalaz = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.scalenie');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.scalenie };
      odswiez();
      return;
    }
    if (!wynik.wynik.merged) {
      const konflikty = wynik.wynik.conflicts ?? [];
      scalenie = { zrodlo: zrodlo.id, cel: celId, konflikty, wybory: new Map() };
      komunikat = { rodzaj: 'blad', tekst: komunikatKonfliktow(konflikty.length) };
      odswiez();
      return;
    }
    scalenie = null;
    komunikat = { rodzaj: 'sukces', tekst: komunikatScalenia(zrodlo.name) };
    const dokument = idDokumentu();
    stanWersji = null;
    if (dokument !== null) await wczytajGalezie(dokument);
    else odswiez();
  }

  // ── Czynności dziennika ────────────────────────────────────────────────

  /* Bilans niesie pominięcia z powodem — przemilczenie ich zostawiłoby
     Operatora w przekonaniu, że cofnięcie objęło całość zmiany. */
  function zdanieBilansu(bilans: StudioActionBalance, podstawowe: string): string {
    const powody = (bilans.skipped ?? []).map((pominiete) => pominiete.reason).filter((powod) => powod !== '');
    const zdanie = bilans.note ?? podstawowe;
    if (powody.length === 0) return zdanie;
    return `${zdanie} ${tresciRepo.dziennik.pominieto} ${powody.join('; ')}`;
  }

  async function cofnijCzynnosc(dokument: string, wpis: StudioDocumentAction): Promise<void> {
    czynnoscWBiegu = wpis.id;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioJournalRevert, {
      documentId: dokument,
      actionIds: [wpis.id],
    });
    if (zdjete) return;
    czynnoscWBiegu = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.cofniecie');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.cofniecie };
      odswiez();
      return;
    }
    const cofniete = wynik.wynik.reverted.length;
    komunikat = {
      rodzaj: cofniete === 0 ? 'blad' : 'sukces',
      tekst: zdanieBilansu(wynik.wynik.balance, komunikatCofniecia(cofniete)),
    };
    stanWersji = null;
    await wczytajDziennik(dokument);
  }

  async function ponowCzynnosc(dokument: string, wpis: StudioDocumentAction): Promise<void> {
    czynnoscWBiegu = wpis.id;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioJournalRedo, {
      documentId: dokument,
      actionIds: [wpis.id],
    });
    if (zdjete) return;
    czynnoscWBiegu = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.ponowienie');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.ponowienie };
      odswiez();
      return;
    }
    const ponowione = wynik.wynik.redone.length;
    komunikat = {
      rodzaj: ponowione === 0 ? 'blad' : 'sukces',
      tekst: zdanieBilansu(wynik.wynik.balance, komunikatPonowienia(ponowione)),
    };
    stanWersji = null;
    await wczytajDziennik(dokument);
  }

  // ── Czynności na kopiach zapasowych ────────────────────────────────────

  async function zalozKopie(dokument: string): Promise<void> {
    zakladanieKopii = true;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBackupCreate, {
      documentId: dokument,
      reason: StudioBackupReason.Manual,
    });
    if (zdjete) return;
    zakladanieKopii = false;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.kopiaZalozenie');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.kopiaZalozenie };
      odswiez();
      return;
    }
    /* Kopia nieudana wraca jako wynik pozytywny z `succeeded` na nie —
       przemilczenie tego zostawiłoby Operatora w przekonaniu, że kopia stoi. */
    komunikat = wynik.wynik.backup.succeeded
      ? { rodzaj: 'sukces', tekst: tresciRepo.komunikat.kopiaZalozona }
      : { rodzaj: 'blad', tekst: tresciRepo.komunikat.kopiaNieudana };
    await wczytajKopie(dokument);
  }

  async function przywrocKopie(dokument: string, kopia: StudioDocumentBackup, doNowego: boolean): Promise<void> {
    kopiaWBiegu = kopia.id;
    odswiez();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioBackupRestore, {
      backupId: kopia.id,
      asNewDocument: doNowego ? true : undefined,
      windowId: zaleznosci.idOkna ?? undefined,
    });
    if (zdjete) return;
    kopiaWBiegu = null;
    if (!wynik.udany || wynik.wynik === undefined) {
      zaloguj(wynik.blad, 'repo.kopiaPrzywrocenie');
      komunikat = { rodzaj: 'blad', tekst: tresciRepo.odmowa.kopiaPrzywrocenie };
      odswiez();
      return;
    }
    /* Nazwa dokumentu przywroconego jest w kontrakcie nieobowiazkowa; jej brak
       nazywa sie wprost, zamiast wstawiac w cudzyslow puste miejsce. */
    const nazwaNowego = wynik.wynik.document.title;
    if (!doNowego) komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.kopiaPrzywrocona };
    else if (nazwaNowego === undefined || nazwaNowego === '')
      komunikat = { rodzaj: 'sukces', tekst: tresciRepo.komunikat.kopiaPrzywroconaNowyBezNazwy };
    else komunikat = { rodzaj: 'sukces', tekst: komunikatKopiiNowyDokument(nazwaNowego) };
    stanWersji = null;
    await wczytajKopie(dokument);
  }

  // ── Widoki wspólne ─────────────────────────────────────────────────────

  function widokLadowanie(opis: string): HTMLElement {
    return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
      el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
      opis,
    ]);
  }

  function widokOdmowa(tytul: string, blad?: ErrorInfo): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('div', { klasa: 'dn-alert-tresc' }, [el('b', { tekst: tytul }), el('span', { tekst: opisOdmowy(blad, 'repo') })]),
    ]);
  }

  function widokKomunikat(k: NonNullable<StanKomunikatu>): HTMLElement {
    const udany = k.rodzaj === 'sukces';
    return el(
      'div',
      {
        klasa: `dn-alert dn-alert--wstega ${udany ? 'dn-alert--sukces' : 'dn-alert--blad'}`,
        role: udany ? 'status' : 'alert',
      },
      [el('span', { klasa: 'dn-alert-tresc', tekst: k.tekst })],
    );
  }

  function widokPusto(tytul: string, opis: string): HTMLElement {
    return el('div', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty' }, [
      znak('pusto'),
      el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tytul }),
      el('span', { klasa: 'dn-pusty-stan-opis', tekst: opis }),
    ]);
  }

  function pigulka(etykieta: string, aktywna: boolean, przelacz: () => void): HTMLElement {
    const guzik = el('button', {
      klasa: aktywna ? 'dn-btn dn-btn--zarys dn-btn--sm' : 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      'aria-pressed': aktywna ? 'true' : 'false',
      tekst: etykieta,
    });
    guzik.addEventListener('click', przelacz);
    return guzik;
  }

  function przelacznikKart(): HTMLElement {
    function karteczka(wartosc: Karta, etykieta: string): HTMLElement {
      return pigulka(etykieta, karta === wartosc, () => {
        if (karta === wartosc) return;
        karta = wartosc;
        komunikat = null;
        formularz = null;
        odswiez();
        wczytajBiezacy(false);
      });
    }
    return el('div', { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': tresciRepo.karty.etykietaPrzelacznika }, [
      karteczka('wersje', tresciRepo.karty.wersje),
      karteczka('galezie', tresciRepo.karty.galezie),
      karteczka('dziennik', tresciRepo.karty.dziennik),
      karteczka('kopie', tresciRepo.karty.kopie),
    ]);
  }

  // ── Widok wersji ───────────────────────────────────────────────────────

  function filtrWersji(): HTMLElement {
    function przelaczSzereg(wartosc: WyborSzeregu): void {
      if (szereg === wartosc) return;
      szereg = wartosc;
      const dokument = idDokumentu();
      if (dokument !== null) void wczytajWersje(dokument);
    }
    const pigulki = el('div', { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': tresciRepo.filtr.szereg }, [
      pigulka(tresciRepo.filtr.wszystkie, szereg === 'wszystkie', () => przelaczSzereg('wszystkie')),
      pigulka(tresciRepo.filtr.operator, szereg === StudioVersionSeries.Operator, () => przelaczSzereg(StudioVersionSeries.Operator)),
      pigulka(tresciRepo.filtr.autozapis, szereg === StudioVersionSeries.Autosave, () => przelaczSzereg(StudioVersionSeries.Autosave)),
      pigulka(tresciRepo.filtr.tylkoKluczowe, tylkoKluczowe, () => {
        tylkoKluczowe = !tylkoKluczowe;
        const dokument = idDokumentu();
        if (dokument !== null) void wczytajWersje(dokument);
      }),
    ]);
    return pigulki;
  }

  function formularzGalaz(dokument: string, wersja: StudioVersion, busy: boolean): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      required: true,
      placeholder: tresciRepo.formularzGalaz.zastepcza,
      'aria-label': tresciRepo.formularzGalaz.etykieta,
    }) as HTMLInputElement;

    const przyciskZaloz = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'submit',
      disabled: busy,
      tekst: busy ? tresciRepo.formularzGalaz.zakladanie : tresciRepo.formularzGalaz.zaloz,
    });
    const przyciskAnuluj = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: tresciRepo.akcje.anuluj,
    });
    przyciskAnuluj.addEventListener('click', () => {
      formularz = null;
      odswiez();
    });

    const formularzWezel = el('form', { klasa: 'dn-pole' }, [
      el('label', { klasa: 'dn-pole-etykieta', tekst: tresciRepo.formularzGalaz.etykieta }),
      pole,
      el('div', { klasa: 'dn-pas-dzialan' }, [przyciskZaloz, przyciskAnuluj]),
    ]);
    formularzWezel.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      const nazwa = pole.value.trim();
      if (nazwa === '') return;
      void zalozGalaz(dokument, wersja, nazwa);
    });
    return formularzWezel;
  }

  function formularzEtykieta(wersja: StudioVersion, busy: boolean): HTMLElement {
    const pole = el('input', {
      klasa: 'dn-pole-kontrolka',
      type: 'text',
      'aria-label': tresciRepo.formularzEtykieta.etykieta,
    }) as HTMLInputElement;
    pole.value = wersja.label ?? '';

    const zaznacz = el('input', { type: 'checkbox', id: `repo-kluczowa-${wersja.id}` }) as HTMLInputElement;
    zaznacz.checked = wersja.milestone === true;
    const zaznaczEtykieta = el('label', { for: `repo-kluczowa-${wersja.id}` }, [
      zaznacz,
      ` ${tresciRepo.formularzEtykieta.kluczowa}`,
    ]);

    const przyciskZapisz = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'submit',
      disabled: busy,
      tekst: busy ? tresciRepo.formularzEtykieta.zapisywanie : tresciRepo.formularzEtykieta.zapisz,
    });
    const przyciskOdwolanie = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      disabled: busy,
      tekst:
        akcjaWiersza?.wersjaId === wersja.id && akcjaWiersza.akcja === 'odwolanie'
          ? tresciRepo.formularzEtykieta.odwolanieWBiegu
          : tresciRepo.formularzEtykieta.odwolanie,
    });
    przyciskOdwolanie.addEventListener('click', () => void utworzOdwolanie(wersja));
    const przyciskAnuluj = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: tresciRepo.akcje.anuluj,
    });
    przyciskAnuluj.addEventListener('click', () => {
      formularz = null;
      odswiez();
    });

    const formularzWezel = el('form', { klasa: 'dn-pole' }, [
      el('label', { klasa: 'dn-pole-etykieta', tekst: tresciRepo.formularzEtykieta.etykieta }),
      pole,
      zaznaczEtykieta,
      el('div', { klasa: 'dn-pas-dzialan' }, [przyciskZapisz, przyciskOdwolanie, przyciskAnuluj]),
    ]);
    formularzWezel.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      void ustawEtykiete(wersja, pole.value.trim(), zaznacz.checked);
    });
    return formularzWezel;
  }

  function wiersz(dokument: string, wersja: StudioVersion, indeks: number, liczbaWersji: number, zalozycielska?: string): HTMLElement {
    const biezaca = indeks === 0;
    const busyTegoWiersza = akcjaWiersza?.wersjaId === wersja.id;

    const naglowekWiersza = el('div', { klasa: 'st-panel-wiersz' }, [
      el('span', {
        klasa: `dn-kropka ${biezaca ? 'dn-kropka--sukces' : 'dn-kropka--neutralna'}`,
        'aria-hidden': 'true',
      }),
      el('b', { tekst: podpisWersji(liczbaWersji - indeks) }),
      el('span', { klasa: 'dn-meta', tekst: `${formatCzas(wersja.createdAt)} · ${opisAutora(wersja.author)}` }),
    ]);

    const dzieci: HTMLElement[] = [naglowekWiersza];

    if (biezaca) {
      dzieci.push(el('div', { klasa: 'dn-nota', tekst: tresciRepo.akcje.biezaca }));
    } else if (wersja.summary !== undefined && wersja.summary !== '') {
      dzieci.push(el('div', { klasa: 'dn-nota', tekst: `„${wersja.summary}”` }));
    }

    const plakietki: HTMLElement[] = [];
    if (wersja.milestone === true || (wersja.label !== undefined && wersja.label !== '')) {
      const tekstPlakietki =
        wersja.milestone === true ? `★ ${wersja.label ?? tresciRepo.wiersz.kluczowa}` : (wersja.label as string);
      plakietki.push(
        el('span', {
          klasa: `dn-plakietka${wersja.milestone === true ? ' dn-plakietka--sukces' : ''}`,
          tekst: tekstPlakietki,
        }),
      );
    }
    if (zalozycielska !== undefined && zalozycielska === wersja.id) {
      plakietki.push(el('span', { klasa: 'dn-plakietka dn-plakietka--informacja', tekst: tresciRepo.wiersz.zalozycielska }));
    }
    if (plakietki.length > 0) dzieci.push(el('div', {}, plakietki));

    const przyciskPrzywroc = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: biezaca || busyTegoWiersza,
      tekst:
        akcjaWiersza?.wersjaId === wersja.id && akcjaWiersza.akcja === 'przywroc'
          ? tresciRepo.akcje.przywracanie
          : tresciRepo.akcje.przywroc,
    });
    przyciskPrzywroc.addEventListener('click', () => void przywroc(dokument, wersja));

    const przyciskRozgalez = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: busyTegoWiersza,
      tekst: tresciRepo.akcje.rozgalez,
    });
    przyciskRozgalez.addEventListener('click', () => {
      formularz =
        formularz?.wersjaId === wersja.id && formularz.akcja === 'rozgalez'
          ? null
          : { akcja: 'rozgalez', wersjaId: wersja.id };
      odswiez();
    });

    const przyciskWiecej = el('button', {
      klasa: 'dn-btn-ikona dn-btn-ikona--sm',
      type: 'button',
      'aria-label': tresciRepo.akcje.wiecej,
      disabled: busyTegoWiersza,
      tekst: '⋯',
    });
    przyciskWiecej.addEventListener('click', () => {
      formularz =
        formularz?.wersjaId === wersja.id && formularz.akcja === 'etykieta'
          ? null
          : { akcja: 'etykieta', wersjaId: wersja.id };
      odswiez();
    });

    dzieci.push(el('div', { klasa: 'dn-wersja-akcje' }, [przyciskPrzywroc, przyciskRozgalez, przyciskWiecej]));

    if (formularz !== null && formularz.wersjaId === wersja.id) {
      dzieci.push(
        formularz.akcja === 'rozgalez'
          ? formularzGalaz(dokument, wersja, busyTegoWiersza)
          : formularzEtykieta(wersja, busyTegoWiersza),
      );
    }

    return el('div', { klasa: 'dn-wersja' }, dzieci);
  }

  function stopkaWersji(dokument: string): HTMLElement {
    const przyciskEksport = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm st-pole-rozciagniete',
      type: 'button',
      disabled: eksportTrwa,
      tekst: eksportTrwa ? tresciRepo.stopka.eksportowanie : tresciRepo.stopka.eksportuj,
    });
    przyciskEksport.addEventListener('click', () => void eksportuj(dokument));

    const przyciskPowrot = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm st-pole-rozciagniete',
      type: 'button',
      disabled: powrotTrwa,
      tekst: powrotTrwa ? tresciRepo.stopka.zalozycielskaWBiegu : tresciRepo.stopka.zalozycielska,
    });
    przyciskPowrot.addEventListener('click', () => void wrocDoZalozycielskiej(dokument));

    return el('div', { klasa: 'dn-wersja-stopka' }, [przyciskEksport, przyciskPowrot]);
  }

  function widokWersji(dokument: string): HTMLElement[] {
    const elementy: HTMLElement[] = [filtrWersji()];
    if (stanWersji === null || stanWersji.rodzaj === 'wczytywanie') {
      elementy.push(widokLadowanie(tresciRepo.stany.wczytywanie));
      return elementy;
    }
    if (stanWersji.rodzaj === 'odmowa') {
      elementy.push(widokOdmowa(tresciRepo.odmowa.lista, stanWersji.blad));
      return elementy;
    }
    if (stanWersji.operatorow !== undefined && stanWersji.autozapisow !== undefined) {
      elementy.push(el('div', { klasa: 'dn-meta', tekst: opisSzeregow(stanWersji.operatorow, stanWersji.autozapisow) }));
    }
    if (stanWersji.wersje.length === 0) {
      const zawezone = szereg !== 'wszystkie' || tylkoKluczowe;
      elementy.push(
        zawezone
          ? widokPusto(tresciRepo.pusto.filtrTytul, tresciRepo.pusto.filtrOpis)
          : widokPusto(tresciRepo.pusto.tytul, tresciRepo.pusto.opis),
      );
      return elementy;
    }
    const zalozycielska = stanWersji.zalozycielska;
    const wersje = stanWersji.wersje;
    wersje.forEach((w, indeks) => elementy.push(wiersz(dokument, w, indeks, wersje.length, zalozycielska)));
    elementy.push(el('div', { klasa: 'dn-nota', tekst: tresciRepo.niegotowe.podgladPorownanie }));
    elementy.push(stopkaWersji(dokument));
    return elementy;
  }

  // ── Widok gałęzi ───────────────────────────────────────────────────────

  function widokKonfliktow(zrodlo: StudioBranch, biezace: NonNullable<StanScalania>): HTMLElement {
    const dzieci: Dziecko[] = [el('div', { klasa: 'pt-etykieta', tekst: tresciRepo.galezie.konfliktyTytul })];
    for (const konflikt of biezace.konflikty) {
      const wybrana = biezace.wybory.get(konflikt.index);
      function strona(wartosc: StudioMergeSide, etykieta: string): HTMLElement {
        return pigulka(etykieta, wybrana === wartosc, () => {
          biezace.wybory.set(konflikt.index, wartosc);
          odswiez();
        });
      }
      dzieci.push(
        el('div', { klasa: 'dn-wersja' }, [
          el('div', { klasa: 'st-panel-wiersz' }, [el('b', { tekst: opisKonfliktu(konflikt.index) })]),
          el('div', { klasa: 'dn-nota', tekst: konflikt.source }),
          el('div', { klasa: 'dn-nota', tekst: konflikt.target }),
          el('div', { klasa: 'dn-wersja-akcje' }, [
            strona(StudioMergeSide.Scalana, tresciRepo.galezie.stronaScalana),
            strona(StudioMergeSide.Docelowa, tresciRepo.galezie.stronaDocelowa),
            strona(StudioMergeSide.Obie, tresciRepo.galezie.stronaObie),
          ]),
        ]),
      );
    }
    dzieci.push(el('div', { klasa: 'dn-nota', tekst: tresciRepo.niegotowe.trescWlasnaKonfliktu }));

    const komplet = biezace.konflikty.every((k) => biezace.wybory.has(k.index));
    const przyciskScal = el('button', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      type: 'button',
      disabled: !komplet || scalanaGalaz !== null,
      tekst: tresciRepo.galezie.scalPonownie,
    });
    przyciskScal.addEventListener('click', () => {
      const rozstrzygniecia: StudioMergeResolution[] = biezace.konflikty.map((k) => ({
        index: k.index,
        side: biezace.wybory.get(k.index) ?? StudioMergeSide.Docelowa,
      }));
      void scal(zrodlo, biezace.cel, rozstrzygniecia);
    });
    const przyciskPorzuc = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      tekst: tresciRepo.galezie.porzucKonflikty,
    });
    przyciskPorzuc.addEventListener('click', () => {
      scalenie = null;
      komunikat = null;
      odswiez();
    });
    dzieci.push(el('div', { klasa: 'dn-pas-dzialan' }, [przyciskScal, przyciskPorzuc]));
    return el('div', { klasa: 'st-panel-lista' }, dzieci);
  }

  function wierszGalezi(galaz: StudioBranch, wszystkie: StudioBranch[]): HTMLElement {
    const dzieci: Dziecko[] = [
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', {
          klasa: `dn-kropka ${galaz.merged ? 'dn-kropka--neutralna' : 'dn-kropka--sukces'}`,
          'aria-hidden': 'true',
        }),
        el('b', { tekst: galaz.name }),
        el('span', { klasa: 'dn-meta', tekst: formatCzas(galaz.createdAt) }),
      ]),
      el('div', { klasa: 'dn-nota', tekst: podpisPunktuStartowego(galaz.fromVersionId) }),
    ];
    if (galaz.merged) {
      dzieci.push(el('div', {}, [el('span', { klasa: 'dn-plakietka', tekst: tresciRepo.galezie.scalona })]));
    }

    const cele = wszystkie.filter((inna) => inna.id !== galaz.id);
    if (cele.length === 0) {
      dzieci.push(el('div', { klasa: 'dn-nota', tekst: tresciRepo.niegotowe.brakCeluScalania }));
      return el('div', { klasa: 'dn-wersja' }, dzieci);
    }

    const wybor = el('select', { klasa: 'dn-pole-kontrolka', 'aria-label': tresciRepo.galezie.docelowa },
      cele.map((inna) => el('option', { value: inna.id, tekst: inna.name })),
    ) as HTMLSelectElement;

    const przyciskScal = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: galaz.merged || scalanaGalaz !== null,
      tekst: scalanaGalaz === galaz.id ? tresciRepo.galezie.scalanie : tresciRepo.galezie.scal,
    });
    przyciskScal.addEventListener('click', () => void scal(galaz, wybor.value));

    dzieci.push(el('div', { klasa: 'dn-wersja-akcje' }, [wybor, przyciskScal]));
    if (scalenie !== null && scalenie.zrodlo === galaz.id) dzieci.push(widokKonfliktow(galaz, scalenie));
    return el('div', { klasa: 'dn-wersja' }, dzieci);
  }

  function widokGalezi(): HTMLElement[] {
    if (stanGalezi === null || stanGalezi.rodzaj === 'wczytywanie') return [widokLadowanie(tresciRepo.stany.wczytywanieGalezi)];
    if (stanGalezi.rodzaj === 'odmowa') return [widokOdmowa(tresciRepo.odmowa.galezie, stanGalezi.blad)];
    if (stanGalezi.galezie.length === 0) return [widokPusto(tresciRepo.pusto.galezieTytul, tresciRepo.pusto.galezieOpis)];
    const wszystkie = stanGalezi.galezie;
    return wszystkie.map((galaz) => wierszGalezi(galaz, wszystkie));
  }

  // ── Widok dziennika ────────────────────────────────────────────────────

  function filtrDziennika(): HTMLElement[] {
    function przeladuj(): void {
      const dokument = idDokumentu();
      if (dokument !== null) void wczytajDziennik(dokument);
    }
    const poAutorze = el('div', { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': tresciRepo.filtr.autor }, [
      pigulka(tresciRepo.filtr.wszystkie, autorWpisow === 'wszyscy', () => {
        autorWpisow = 'wszyscy';
        przeladuj();
      }),
      pigulka(tresciRepo.wiersz.autorUzytkownik, autorWpisow === StudioAuthor.Uzytkownik, () => {
        autorWpisow = StudioAuthor.Uzytkownik;
        przeladuj();
      }),
      pigulka(tresciRepo.wiersz.autorModel, autorWpisow === StudioAuthor.Model, () => {
        autorWpisow = StudioAuthor.Model;
        przeladuj();
      }),
    ]);
    const poStanie = el('div', { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': tresciRepo.filtr.stan }, [
      pigulka(tresciRepo.filtr.wszystkie, stanWpisow === 'wszystkie', () => {
        stanWpisow = 'wszystkie';
        przeladuj();
      }),
      pigulka(tresciRepo.filtr.stanCzynne, stanWpisow === StudioActionState.Active, () => {
        stanWpisow = StudioActionState.Active;
        przeladuj();
      }),
      pigulka(tresciRepo.filtr.stanCofniete, stanWpisow === StudioActionState.Reverted, () => {
        stanWpisow = StudioActionState.Reverted;
        przeladuj();
      }),
    ]);
    return [poAutorze, poStanie];
  }

  function wierszDziennika(dokument: string, wpis: StudioDocumentAction): HTMLElement {
    const cofnieta = wpis.state === StudioActionState.Reverted;
    const busy = czynnoscWBiegu === wpis.id;
    const zablokowana = (wpis.blockedBy ?? []).length > 0;

    const dzieci: Dziecko[] = [
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', {
          klasa: `dn-kropka ${cofnieta ? 'dn-kropka--neutralna' : 'dn-kropka--sukces'}`,
          'aria-hidden': 'true',
        }),
        el('b', { tekst: wpis.description }),
        el('span', {
          klasa: 'dn-meta',
          tekst: `${formatCzas(wpis.createdAt)} · ${opisAutora(wpis.author)} · ${opisRodzaju(wpis.kind)}`,
        }),
      ]),
    ];

    const plakietki: HTMLElement[] = [];
    if (cofnieta) plakietki.push(el('span', { klasa: 'dn-plakietka', tekst: tresciRepo.dziennik.cofnieta }));
    if (zablokowana) {
      plakietki.push(el('span', { klasa: 'dn-plakietka dn-plakietka--ostrzezenie', tekst: tresciRepo.dziennik.zaleznosc }));
    }
    if (plakietki.length > 0) dzieci.push(el('div', {}, plakietki));

    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: busy,
      tekst: cofnieta
        ? busy
          ? tresciRepo.dziennik.ponawianie
          : tresciRepo.dziennik.ponow
        : busy
          ? tresciRepo.dziennik.cofanie
          : tresciRepo.dziennik.cofnij,
    });
    przycisk.addEventListener('click', () => {
      if (cofnieta) void ponowCzynnosc(dokument, wpis);
      else void cofnijCzynnosc(dokument, wpis);
    });
    dzieci.push(el('div', { klasa: 'dn-wersja-akcje' }, [przycisk]));
    return el('div', { klasa: 'dn-wersja' }, dzieci);
  }

  function widokDziennika(dokument: string): HTMLElement[] {
    const elementy: HTMLElement[] = [...filtrDziennika()];
    if (stanDziennika === null || stanDziennika.rodzaj === 'wczytywanie') {
      elementy.push(widokLadowanie(tresciRepo.stany.wczytywanieDziennika));
      return elementy;
    }
    if (stanDziennika.rodzaj === 'odmowa') {
      elementy.push(widokOdmowa(tresciRepo.odmowa.dziennik, stanDziennika.blad));
      return elementy;
    }
    if (stanDziennika.wpisy.length === 0) {
      const zawezone = autorWpisow !== 'wszyscy' || stanWpisow !== 'wszystkie';
      elementy.push(
        zawezone
          ? widokPusto(tresciRepo.pusto.filtrTytul, tresciRepo.pusto.filtrOpis)
          : widokPusto(tresciRepo.pusto.dziennikTytul, tresciRepo.pusto.dziennikOpis),
      );
      return elementy;
    }
    elementy.push(el('div', { klasa: 'dn-meta', tekst: opisWpisow(stanDziennika.wpisy.length, stanDziennika.wszystkich) }));
    for (const wpis of stanDziennika.wpisy) elementy.push(wierszDziennika(dokument, wpis));
    return elementy;
  }

  // ── Widok kopii zapasowych ─────────────────────────────────────────────

  function wierszKopii(dokument: string, kopia: StudioDocumentBackup): HTMLElement {
    const busy = kopiaWBiegu === kopia.id;
    const dzieci: Dziecko[] = [
      el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', {
          klasa: `dn-kropka ${kopia.succeeded ? 'dn-kropka--sukces' : 'dn-kropka--blad'}`,
          'aria-hidden': 'true',
        }),
        el('b', { tekst: opisPowodu(kopia.reason) }),
        el('span', {
          klasa: 'dn-meta',
          tekst: kopia.bytes === undefined ? formatCzas(kopia.createdAt) : `${formatCzas(kopia.createdAt)} · ${formatBajty(kopia.bytes)}`,
        }),
      ]),
    ];

    const plakietki: HTMLElement[] = [];
    if (kopia.unsavedChanges === true) {
      plakietki.push(el('span', { klasa: 'dn-plakietka dn-plakietka--ostrzezenie', tekst: tresciRepo.kopie.niezapisane }));
    }
    if (!kopia.succeeded) {
      plakietki.push(el('span', { klasa: 'dn-plakietka dn-plakietka--blad', tekst: tresciRepo.kopie.nieudana }));
    }
    if (plakietki.length > 0) dzieci.push(el('div', {}, plakietki));
    if (kopia.failureReason !== undefined && kopia.failureReason !== '') {
      dzieci.push(el('div', { klasa: 'dn-nota', tekst: kopia.failureReason }));
    }

    const przyciskNaMiejsce = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: busy || !kopia.succeeded,
      tekst: busy ? tresciRepo.kopie.przywracanie : tresciRepo.kopie.przywroc,
    });
    przyciskNaMiejsce.addEventListener('click', () => void przywrocKopie(dokument, kopia, false));

    const przyciskNowy = el('button', {
      klasa: 'dn-btn dn-btn--duch dn-btn--sm',
      type: 'button',
      disabled: busy || !kopia.succeeded,
      tekst: tresciRepo.kopie.doNowego,
    });
    przyciskNowy.addEventListener('click', () => void przywrocKopie(dokument, kopia, true));

    dzieci.push(el('div', { klasa: 'dn-wersja-akcje' }, [przyciskNaMiejsce, przyciskNowy]));
    return el('div', { klasa: 'dn-wersja' }, dzieci);
  }

  function widokKopii(dokument: string): HTMLElement[] {
    const przyciskZaloz = el('button', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      type: 'button',
      disabled: zakladanieKopii,
      tekst: zakladanieKopii ? tresciRepo.kopie.zakladanie : tresciRepo.kopie.zaloz,
    });
    przyciskZaloz.addEventListener('click', () => void zalozKopie(dokument));

    const filtr = el('div', { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': tresciRepo.filtr.niezapisane }, [
      pigulka(tresciRepo.filtr.wszystkie, !tylkoNiezapisane, () => {
        if (!tylkoNiezapisane) return;
        tylkoNiezapisane = false;
        void wczytajKopie(dokument);
      }),
      pigulka(tresciRepo.filtr.niezapisane, tylkoNiezapisane, () => {
        if (tylkoNiezapisane) return;
        tylkoNiezapisane = true;
        void wczytajKopie(dokument);
      }),
    ]);

    const elementy: HTMLElement[] = [filtr, el('div', { klasa: 'dn-pas-dzialan' }, [przyciskZaloz])];
    if (stanKopii === null || stanKopii.rodzaj === 'wczytywanie') {
      elementy.push(widokLadowanie(tresciRepo.stany.wczytywanieKopii));
      return elementy;
    }
    if (stanKopii.rodzaj === 'odmowa') {
      elementy.push(widokOdmowa(tresciRepo.odmowa.kopie, stanKopii.blad));
      return elementy;
    }
    elementy.push(el('div', { klasa: 'dn-meta', tekst: opisNiezapisanych(stanKopii.niezapisanych) }));
    if (stanKopii.kopie.length === 0) {
      elementy.push(
        tylkoNiezapisane
          ? widokPusto(tresciRepo.pusto.filtrTytul, tresciRepo.pusto.filtrOpis)
          : widokPusto(tresciRepo.pusto.kopieTytul, tresciRepo.pusto.kopieOpis),
      );
      return elementy;
    }
    for (const kopia of stanKopii.kopie) elementy.push(wierszKopii(dokument, kopia));
    return elementy;
  }

  // ── Złożenie widoku ────────────────────────────────────────────────────

  function zawartoscKarty(dokument: string): HTMLElement[] {
    switch (karta) {
      case 'wersje':
        return widokWersji(dokument);
      case 'galezie':
        return widokGalezi();
      case 'dziennik':
        return widokDziennika(dokument);
      case 'kopie':
        return widokKopii(dokument);
    }
  }

  function zawartosc(): HTMLElement[] {
    if (podstawa.rodzaj === 'brakOkna') return [widokOdmowa(tresciRepo.stany.brakOkna)];
    if (podstawa.rodzaj === 'oczekiwanieDokumentu') return [widokLadowanie(tresciRepo.stany.oczekiwanieDokumentu)];
    const elementy: HTMLElement[] = [przelacznikKart()];
    if (komunikat !== null) elementy.push(widokKomunikat(komunikat));
    elementy.push(...zawartoscKarty(podstawa.id));
    return elementy;
  }

  function odswiez(): void {
    tresc.replaceChildren(...zawartosc());
    const pokazLicznik = karta === 'wersje' && stanWersji !== null && stanWersji.rodzaj === 'gotowe';
    znacznik.hidden = !pokazLicznik;
    znacznik.textContent =
      stanWersji !== null && stanWersji.rodzaj === 'gotowe' && pokazLicznik ? znacznikWersji(stanWersji.wersje.length) : '';
  }

  odswiez();

  const odsubskrybuj: (() => void) | null =
    zaleznosci.idOkna === null
      ? null
      : zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
          if (zdjete || zdarzenie.document.windowId !== zaleznosci.idOkna) return;
          if (zdarzenie.change === ChangeKind.Deleted) {
            podstawa = { rodzaj: 'oczekiwanieDokumentu' };
            stanWersji = null;
            stanGalezi = null;
            stanDziennika = null;
            stanKopii = null;
            komunikat = null;
            odswiez();
            return;
          }
          /* Zmiana dokumentu unieważnia wszystkie cztery widoki, ale wczytany
             zostaje tylko ten, na który Operator patrzy — pozostałe wczytają
             się przy przełączeniu, żeby zdarzenie nie wołało czterech komend. */
          const nowy = zdarzenie.document.id;
          const inny = podstawa.rodzaj !== 'dokument' || podstawa.id !== nowy;
          podstawa = { rodzaj: 'dokument', id: nowy };
          if (inny) {
            stanGalezi = null;
            stanDziennika = null;
            stanKopii = null;
            scalenie = null;
            formularz = null;
          }
          stanWersji = null;
          odswiez();
          wczytajBiezacy(true);
        });

  return {
    zdejmij() {
      zdjete = true;
      odsubskrybuj?.();
      wezel.replaceChildren();
    },
  };
}
