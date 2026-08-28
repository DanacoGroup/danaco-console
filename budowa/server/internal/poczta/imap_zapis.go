// Odpowiedzialność pliku: ZAPIS do skrzynki Operatora — szkic odpowiedzi
// i oznaczenie listu; obie czynności robią dokładnie to, co zlecono, i nic
// ponadto.
package poczta

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap/v2"
)

// ZapiszSzkic odkłada szkic w folderze szkiców i oddaje go w kształcie
// nagłówka; wskazanie poprzedni (niepuste) kasuje poprzednią wersję szkicu.
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
		// nowy szkic już leży w skrzynce.
		_ = k.usunWiadomosc(poprzedni)
	}
	return k.naglowekPo(folder, uid)
}

// Oznacz zmienia znaczniki listu — obsługuje mail.message.flag; wskazanie
// puste (nil) zostawia znacznik NIETKNIĘTY, zgodnie z opisem kontraktu.
func (k *Klient) Oznacz(identyfikator string, nieprzeczytana, wyrozniona *bool) (Naglowek, error) {
	folder, uid, err := rozbierzIdentyfikator(identyfikator)
	if err != nil {
		return Naglowek{}, err
	}
	if _, err := k.imap.Select(folder, nil).Wait(); err != nil {
		return Naglowek{}, fmt.Errorf("skrzynka nie ma folderu %q albo nie daje do niego zapisu: %w", folder, err)
	}

	// Dwie osobne komendy STORE, bo jedna dokłada znaczniki, a druga je
	// zdejmuje.
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

// znacznikiZmiany przekłada żądanie kontraktu na dwie listy znaczników
// IMAP; nieprzeczytana jest odwrotnością znacznika Seen w protokole.
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

// dolozDoFolderu wykonuje APPEND i oddaje UID dołożonej wiadomości; UID
// bywa nieznany, gdy serwer nie wspiera rozszerzenia UIDPLUS.
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

// usunWiadomosc kasuje wiadomość na dobre: znacznik Deleted plus EXPUNGE,
// bo sam znacznik zostawiłby wpis widoczny dla Operatora.
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

// odnajdzFolder pyta serwer o folder o danym przeznaczeniu (Drafts, Sent)
// i schodzi na nazwę domyślną dopiero wtedy, gdy serwer znacznika nie oddał.
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
	// Folder domyślny musi ISTNIEĆ przed dołożeniem; CREATE na folderze
	// zastanym pomijane świadomie.
	_ = k.imap.Create(domyslny, nil).Wait()
	return domyslny, nil
}

// naglowekPo odczytuje wiadomość PO zmianie, żeby odpowiedź opisywała stan
// rzeczywisty, a nie zamierzony; UID zerowy bierze wiadomość najświeższą.
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
