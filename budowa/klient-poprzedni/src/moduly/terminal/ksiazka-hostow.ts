/**
 * Książka hostów okna Session Manager — wpisy połączeń zdalnych wraz z czytaniem i zapisem
 * pliku konfiguracyjnego OpenSSH; jest widokiem na dziennik rdzenia, nie jego zamiennikiem.
 */

/** Jeden wpis książki hostów okna Session Manager, z adresem, portem, folderem i notatką pochodzenia wpisu. */
export interface WpisHosta {
  /** Identyfikator wpisu nadany przez rdzeń; pusty przed pierwszym zapisem w dzienniku. */
  id?: string;
  /** Port połączenia; brak bierze port domyślny protokołu. */
  port?: number;
  /** Nazwa wpisu widoczna w wykazie i w nazwie karty. */
  nazwa: string;
  /** Adres przekazywany programowi ssh: `użytkownik@host` albo alias konfiguracji serwera. */
  cel: string;
  /** Folder porządkujący wykaz; pusty znaczy „bez folderu”. */
  grupa: string;
  /** Katalog roboczy karty zdalnej; pusty znaczy „katalog własny okna”. */
  katalog: string;
  /** Notatka Operatora albo ślad pochodzenia wpisu. */
  notatka: string;
}

export interface KsiazkaHostow {
  /** Wpisy w kolejności dodania. */
  wpisy(): readonly WpisHosta[];
  /** Wpisy pogrupowane po folderze, foldery w porządku alfabetycznym. */
  grupy(): Array<[string, WpisHosta[]]>;
  /** Dokłada wpis; nazwa jest kluczem, więc wpis o nazwie już zajętej podmienia poprzedni. */
  dodaj(wpis: WpisHosta): boolean;
  /** Usuwa wpis o tej nazwie; fałsz znaczy „takiego wpisu nie ma”. */
  usun(nazwa: string): boolean;
  /** Dokłada wiele wpisów naraz; oddaje liczbę wpisów nowych i podmienionych. */
  scal(wpisy: readonly WpisHosta[]): { nowe: number; podmienione: number };
  /** Zastępuje całą zawartość wykazem z rdzenia; prawdą o książce jest dziennik rdzenia. */
  zastap(wpisy: readonly WpisHosta[]): void;
}

export function utworzKsiazkeHostow(): KsiazkaHostow {
  const wpisy: WpisHosta[] = [];

  function miejsce(nazwa: string): number {
    return wpisy.findIndex((wpis) => wpis.nazwa === nazwa);
  }

  return {
    wpisy: () => wpisy,

    grupy() {
      const grupy = new Map<string, WpisHosta[]>();
      for (const wpis of wpisy) {
        const nazwa = wpis.grupa === '' ? 'Bez folderu' : wpis.grupa;
        const pozycje = grupy.get(nazwa) ?? [];
        pozycje.push(wpis);
        grupy.set(nazwa, pozycje);
      }
      return [...grupy.entries()].sort(([a], [b]) => a.localeCompare(b, 'pl'));
    },

    dodaj(wpis) {
      const stare = miejsce(wpis.nazwa);
      if (stare < 0) {
        wpisy.push(wpis);
        return true;
      }
      wpisy.splice(stare, 1, wpis);
      return false;
    },

    usun(nazwa) {
      const stare = miejsce(nazwa);
      if (stare < 0) return false;
      wpisy.splice(stare, 1);
      return true;
    },

    zastap(nowe) {
      wpisy.splice(0, wpisy.length, ...nowe);
    },

    scal(nowe) {
      let policzoneNowe = 0;
      let policzonePodmienione = 0;
      for (const wpis of nowe) {
        const stare = miejsce(wpis.nazwa);
        if (stare < 0) {
          wpisy.push(wpis);
          policzoneNowe += 1;
          continue;
        }
        wpisy.splice(stare, 1, wpis);
        policzonePodmienione += 1;
      }
      return { nowe: policzoneNowe, podmienione: policzonePodmienione };
    },
  };
}

/**
 * Czyta plik konfiguracyjny OpenSSH i składa z niego wpisy książki; czytany jest podzbiór
 * składni, który niesie adres.
 */
export function czytajKonfiguracjeSsh(tresc: string): WpisHosta[] {
  const wpisy: WpisHosta[] = [];
  let biezacy: { alias: string; host: string; uzytkownik: string; port: string } | null = null;

  function domknij(): void {
    if (biezacy === null) return;
    const { alias, host, uzytkownik, port } = biezacy;
    biezacy = null;
    const maszyna = host === '' ? alias : host;
    const cel = uzytkownik === '' ? maszyna : `${uzytkownik}@${maszyna}`;
    const czesci = [`wpis z konfiguracji OpenSSH (alias ${alias})`];
    const numerPortu = Number.parseInt(port, 10);
    const wpis: WpisHosta = {
      nazwa: alias,
      cel,
      grupa: 'Z konfiguracji OpenSSH',
      katalog: '',
      notatka: czesci.join('; '),
    };
    // Port wchodzi do wpisu polem `remotePort`; wartość niepoprawna nie wchodzi wcale.
    if (Number.isFinite(numerPortu) && numerPortu > 0 && numerPortu <= 65535) {
      wpis.port = numerPortu;
    }
    wpisy.push(wpis);
  }

  for (const surowy of tresc.split('\n')) {
    const wiersz = surowy.replace(/#.*$/, '').trim();
    if (wiersz === '') continue;
    const rozbite = wiersz.split(/[\s=]+/);
    const slowo = (rozbite[0] ?? '').toLowerCase();
    const wartosc = rozbite.slice(1).join(' ').trim();
    if (slowo === 'host') {
      domknij();
      // Blok Host bywa wielokrotny; książka bierze pierwszy alias jako nazwę wpisu.
      const alias = wartosc.split(/\s+/)[0] ?? '';
      if (alias === '' || alias.includes('*') || alias.includes('?')) continue;
      biezacy = { alias, host: '', uzytkownik: '', port: '' };
      continue;
    }
    if (biezacy === null) continue;
    if (slowo === 'hostname') biezacy.host = wartosc;
    else if (slowo === 'user') biezacy.uzytkownik = wartosc;
    else if (slowo === 'port') biezacy.port = wartosc;
  }
  domknij();
  return wpisy;
}

/**
 * Składa treść pliku konfiguracyjnego OpenSSH z wpisów książki; zapisywane jest wyłącznie to,
 * co książka naprawdę wie.
 */
export function zapiszKonfiguracjeSsh(wpisy: readonly WpisHosta[]): string {
  const wiersze: string[] = [
    '# Książka hostów okna Session Manager — Danaco Console, moduł Terminal.',
    '# Plik niesie wyłącznie alias, maszynę i użytkownika. Port, klucz i ProxyJump',
    '# pozostają w konfiguracji maszyny, na której działa program ssh.',
  ];
  for (const wpis of wpisy) {
    const podzial = wpis.cel.split('@');
    const uzytkownik = podzial.length > 1 ? (podzial[0] ?? '') : '';
    const maszyna = podzial.length > 1 ? podzial.slice(1).join('@') : wpis.cel;
    wiersze.push('', `Host ${wpis.nazwa}`, `  HostName ${maszyna}`);
    if (uzytkownik !== '') wiersze.push(`  User ${uzytkownik}`);
    if (wpis.port !== undefined) wiersze.push(`  Port ${wpis.port}`);
    if (wpis.notatka !== '') wiersze.push(`  # ${wpis.notatka.replace(/\n/g, ' ')}`);
  }
  return `${wiersze.join('\n')}\n`;
}
