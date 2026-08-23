// Odpowiedzialność pliku: ZAPIS do skrzynki Operatora — szkic odpowiedzi
// i oznaczenie listu. Obie czynności zmieniają CUDZĄ skrzynkę, więc obie
// robią dokładnie to, co zlecono, i nic ponadto.
//
// Szkic zapisujemy na serwerze, nie u siebie. Szkic trzymany w bazie rdzenia
// byłby szkicem, którego Operator nie zobaczy w swoim kliencie poczty — a cały
// sens tego kroku polega na tym, żeby zobaczył go tam, gdzie zawsze, zanim
// cokolwiek wyjdzie w świat.
//
// IMAP nie zna poprawiania w miejscu. Zmiana szkicu to APPEND nowej wersji
// i skasowanie starej — tak robi każdy klient poczty, bo protokół nie daje nic
// innego. Kolejność jest tu zamierzona: NAJPIERW dokładamy nowy, POTEM
// kasujemy stary. Odwrotnie awaria w połowie zostawiłaby Operatora bez obu
// wersji odpowiedzi, którą model dla niego napisał.
package poczta

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
)

// ZapiszSzkic odkłada szkic w folderze szkiców i oddaje go w kształcie
// nagłówka. Wskazanie `poprzedni` (niepuste) kasuje poprzednią wersję szkicu.
func (k *Klient) ZapiszSzkic(w Wychodzacy, poprzedni string) (Naglowek, error) {
	dokument, err := zlozWychodzacy(w)
	if err != nil {
		return Naglowek{}, err
	}
	folder, err := k.odnajdzFolder(imap.MailboxAttrDrafts, FolderSzkicow)
	if err != nil {
		return Naglowek{}, err
	}

	uid, err := k.dolozDoFolderu(folder, dokument, []imap.Flag{imap.FlagDraft, imap.FlagSeen})
	if err != nil {
		return Naglowek{}, err
	}
	if poprzedni != "" {
		// Niepowodzenie kasowania POPRZEDNIEJ wersji nie jest odmową komendy:
		// nowy szkic już leży w skrzynce, więc praca została wykonana. Operator
		// zobaczy wtedy dwie wersje — co jest stanem gorszym niż jedna, ale
		// nieporównanie lepszym niż odmowa po udanym zapisie.
		_ = k.usunWiadomosc(poprzedni)
	}
	return k.naglowekPo(folder, uid)
}

// Oznacz zmienia znaczniki listu — obsługuje `mail.message.flag`. Wskazanie
// puste (nil) zostawia znacznik NIETKNIĘTY: kontrakt opisuje brak zmiany
// brakiem pola, więc ustawienie go na wartość domyślną byłoby zmianą, której
// nikt nie zlecił.
func (k *Klient) Oznacz(identyfikator string, nieprzeczytana, wyrozniona *bool) (Naglowek, error) {
	folder, uid, err := rozbierzIdentyfikator(identyfikator)
	if err != nil {
		return Naglowek{}, err
	}
	if _, err := k.imap.Select(folder, nil).Wait(); err != nil {
		return Naglowek{}, fmt.Errorf("skrzynka nie ma folderu %q albo nie daje do niego zapisu: %w", folder, err)
	}

	// Dwie osobne komendy STORE, bo jedna dokłada znaczniki, a druga je zdejmuje
	// — a żądanie ma prawo zrobić naraz jedno i drugie (np. „przeczytana
	// i wyróżniona" z listu nieprzeczytanego i niewyróżnionego).
	dokladane, zdejmowane := znacznikiZmiany(nieprzeczytana, wyrozniona)
	for operacja, znaczniki := range map[imap.StoreFlagsOp][]imap.Flag{
		imap.StoreFlagsAdd: dokladane,
		imap.StoreFlagsDel: zdejmowane,
	} {
		if len(znaczniki) == 0 {
			continue
		}
		polecenie := k.imap.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
			Op: operacja, Silent: true, Flags: znaczniki,
		}, nil)
		if err := polecenie.Close(); err != nil {
			return Naglowek{}, fmt.Errorf("serwer poczty odmówił zmiany oznaczeń wiadomości %s: %w",
				identyfikator, err)
		}
	}
	return k.naglowekPo(folder, uid)
}

// znacznikiZmiany przekłada żądanie kontraktu na dwie listy znaczników IMAP.
//
// Nieprzeczytana jest odwrotnością \Seen — IMAP zna wyłącznie „przeczytana”.
// Żądanie `unread: true` ZDEJMUJE więc znacznik, a `unread: false` go dokłada;
// pomylenie tych dwóch kierunków oznaczałoby, że rdzeń oznacza listy dokładnie
// na odwrót, niż prosi Operator.
func znacznikiZmiany(nieprzeczytana, wyrozniona *bool) (dokladane, zdejmowane []imap.Flag) {
	if nieprzeczytana != nil {
		if *nieprzeczytana {
			zdejmowane = append(zdejmowane, imap.FlagSeen)
		} else {
			dokladane = append(dokladane, imap.FlagSeen)
		}
	}
	if wyrozniona != nil {
		if *wyrozniona {
			dokladane = append(dokladane, imap.FlagFlagged)
		} else {
			zdejmowane = append(zdejmowane, imap.FlagFlagged)
		}
	}
	return dokladane, zdejmowane
}

// dolozDoFolderu wykonuje APPEND i oddaje UID dołożonej wiadomości.
//
// UID bywa nieznany i to nie jest błąd. Rozszerzenie UIDPLUS (którym serwer
// odpowiada `APPENDUID`) jest opcjonalne; serwer, który go nie ma, kończy
// APPEND powodzeniem bez podania UID-u. Zwracamy wtedy zero, a wołający
// odnajduje wiadomość osobnym szukaniem — zamiast odmawiać czynności, która
// się udała.
func (k *Klient) dolozDoFolderu(folder string, dokument []byte, znaczniki []imap.Flag) (imap.UID, error) {
	polecenie := k.imap.Append(folder, int64(len(dokument)), &imap.AppendOptions{
		Flags: znaczniki,
		Time:  time.Now(),
	})
	if _, err := polecenie.Write(dokument); err != nil {
		return 0, fmt.Errorf("serwer poczty przerwał zapis do folderu %q: %w", folder, err)
	}
	if err := polecenie.Close(); err != nil {
		return 0, fmt.Errorf("serwer poczty przerwał zapis do folderu %q: %w", folder, err)
	}
	dane, err := polecenie.Wait()
	if err != nil {
		return 0, fmt.Errorf("serwer poczty odmówił zapisu do folderu %q: %w", folder, err)
	}
	if dane == nil {
		return 0, nil
	}
	return dane.UID, nil
}

// usunWiadomosc kasuje wiadomość na dobre: znacznik \Deleted plus EXPUNGE.
// Sam znacznik zostawiłby w folderze szkiców wpis, który Operator dalej widzi.
func (k *Klient) usunWiadomosc(identyfikator string) error {
	folder, uid, err := rozbierzIdentyfikator(identyfikator)
	if err != nil {
		return err
	}
	if _, err := k.imap.Select(folder, nil).Wait(); err != nil {
		return err
	}
	polecenie := k.imap.Store(imap.UIDSetNum(uid), &imap.StoreFlags{
		Op: imap.StoreFlagsAdd, Silent: true, Flags: []imap.Flag{imap.FlagDeleted},
	}, nil)
	if err := polecenie.Close(); err != nil {
		return err
	}
	return k.imap.Expunge().Close()
}

// odnajdzFolder pyta SERWER o folder o danym przeznaczeniu (\Drafts, \Sent)
// i schodzi na nazwę domyślną dopiero wtedy, gdy serwer znacznika nie oddał.
//
// Po co to pytanie. Folder szkiców nazywa się „Drafts" na jednym serwerze,
// „INBOX.Drafts" na drugim i „[Gmail]/Wersje robocze" na trzecim. Zaszycie
// którejkolwiek nazwy oznaczałoby, że szkic ląduje w NOWYM folderze o nazwie
// „Drafts", którego Operator nigdy nie otwiera — czyli że znika.
func (k *Klient) odnajdzFolder(przeznaczenie imap.MailboxAttr, domyslny string) (string, error) {
	wykaz, err := k.imap.List("", "*", &imap.ListOptions{SelectSpecialUse: true}).Collect()
	if err == nil {
		for _, folder := range wykaz {
			for _, cecha := range folder.Attrs {
				if cecha == przeznaczenie {
					return folder.Mailbox, nil
				}
			}
		}
	}
	// Folder domyślny musi ISTNIEĆ, zanim cokolwiek do niego dołożymy. CREATE
	// na folderze zastanym kończy się błędem, który tu pomijamy świadomie —
	// „już jest" to dokładnie ten stan, o który nam chodziło.
	_ = k.imap.Create(domyslny, nil).Wait()
	return domyslny, nil
}

// naglowekPo odczytuje wiadomość PO zmianie, żeby odpowiedź opisywała stan
// rzeczywisty, a nie zamierzony. UID zerowy (serwer bez UIDPLUS)
// znaczy „weź najświeższą wiadomość folderu" — tę, którą właśnie dołożyliśmy.
func (k *Klient) naglowekPo(folder string, uid imap.UID) (Naglowek, error) {
	if uid == 0 {
		wykaz, _, err := k.Wykaz(Zawezenie{Folder: folder, Granica: 1})
		if err != nil {
			return Naglowek{}, err
		}
		if len(wykaz) == 0 {
			return Naglowek{}, fmt.Errorf("folder %q jest pusty tuż po zapisie — "+
				"serwer poczty przyjął wiadomość, ale jej nie pokazuje", folder)
		}
		return wykaz[0], nil
	}
	list, err := k.Pobierz(zlozIdentyfikator(folder, uid), false)
	if err != nil {
		return Naglowek{}, err
	}
	return list.Naglowek, nil
}
