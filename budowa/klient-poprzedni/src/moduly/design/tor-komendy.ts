import { Command } from '../../../../shared/contract';

/**
 * Co moduł wie o torze swoich komend ponad to, co rdzeń mówi sam.
 *
 * Słownik zdań doklejanych do powodu odmowy. Stoi osobno od
 * `etykiety-designu.ts`, bo tamten plik jest słownikiem nazw i braków, a to
 * jest wiedza o drodze żądania.
 *
 * Zdanie nie orzeka o powodzie odmowy — powód niesie `opisOdmowyBledu`
 * z odpowiedzi rdzenia. Tu stoi wyłącznie to, czego rdzeń o sobie nie mówi.
 */
const TORY_KOMEND: Readonly<Record<string, string>> = {
  [Command.DesignAssetGenerate]:
    'Tor komendy: rdzeń sprawdza najpierw RODZAJ wskazanego kanału (pole channelId; bez ' +
    'wskazania bierze pierwszy czynny kanał obrazowy konta), potem wysyła polecenie i czeka ' +
    'na fragment obrazu, odkłada bajty w swoim magazynie pod sumą kontrolną, mierzy format ' +
    'i wymiary z nagłówka utrwalonego pliku, a wiersz zasobu zakłada DOPIERO wtedy. Odmowa ' +
    'na którymkolwiek kroku nie zostawia zasobu bez treści. Kanał tekstowy jest odmawiany ' +
    'przed wysłaniem czegokolwiek: fragment tekstu nie jest obrazem.',
  [Command.DesignAssetList]:
    'Tor komendy: odczyt zawęża się polami windowId, kind, tags, favoriteOnly i limit. ' +
    'Fraza wyszukiwania nie jest polem żądania i zawęża wyłącznie wykaz już wczytany.',
  [Command.DesignBoardUpdate]:
    'Tor komendy: rdzeń przyjmuje CAŁY układ w polu layers naraz — zapis zastępuje układ, ' +
    'a nie dokłada się do niego. Drogą powrotną jest design.board.list.',
  [Command.DesignAssetTagSet]:
    'Tor komendy: zestaw etykiet ZASTĘPUJE poprzedni, także zestawem pustym. Rdzeń odmawia ' +
    'zasobu, którego nie zna, i mówi to wprost.',
  [Command.DesignAssetUpload]:
    'Tor komendy: treść idzie bajtami w polu contentBase64 i wchodzi do magazynu rdzenia pod ' +
    'sumą kontrolną — zasób przestaje zależeć od pliku na dysku Operatora. Pola sourcePath ' +
    'klient nie używa: kazałoby ono otworzyć plik RDZENIOWI, który stoi na innej maszynie.',
  [Command.DesignAssetRemove]:
    'Tor komendy: usunięcie jest natychmiastowe i nieodwracalne — kontrakt nie zna ani ' +
    'archiwizacji, ani przywrócenia. Rdzeń rozgłasza je zdarzeniem design.asset.changed ' +
    'o rodzaju deleted, więc drugie okno zdejmuje zasób samo.',
  [Command.DesignAssetFavoriteSet]:
    'Tor komendy: pole favorite jest stanem DOCELOWYM, nie przestawieniem — żądanie mówi, ' +
    'jak ma być, a nie „zmień na przeciwne". Po nim design.asset.list zawęża po tym polu.',
  [Command.DesignBoardList]:
    'Tor komendy: odczyt zawęża się WYŁĄCZNIE oknem modułu i oddaje kompozycje wraz z ich ' +
    'warstwami. Rdzeń nie obiecuje porządku wykazu — o świeżości rozstrzyga pole updatedAt.',
};

/** Zdanie o torze komendy; komenda spoza obszaru nie dostaje dopisku. */
export function zdanieOTorzeKomendy(komenda: string): string {
  return TORY_KOMEND[komenda] ?? '';
}

/**
 * Powód odmowy podany przez rdzeń wraz ze zdaniem o torze komendy.
 *
 * Kolejność jest stała: najpierw to, co powiedział rdzeń, potem to, co wie
 * okno. Odwrotnie czytelnik dostałby najpierw uogólnienie, a dopiero pod nim
 * rzecz, która go dotyczy.
 */
export function powodZTorem(powod: string, komenda: string): string {
  const tor = zdanieOTorzeKomendy(komenda);
  return tor === '' ? powod : `${powod} ${tor}`;
}
