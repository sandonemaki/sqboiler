package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sqlboiler-project/models"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// searchBooks関数はmain関数の外で定義
func searchBooks(ctx context.Context, db *sql.DB, keyword string) ([]*models.Book, error) {
	return models.Books(
		qm.Where("title LIKE ? OR author LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
		),
		qm.OrderBy("title ASC"),
	).All(ctx, db)
}

func main() {
	// データベース接続
	db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=sqlboiler_db sslmode=disable")
	if err != nil {
		log.Printf("データベース接続エラー: %v\n", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	// エラーチェックを追加
	books, err := models.Books().All(ctx, db)
	if err != nil {
		log.Printf("本の取得エラー: %v\n", err)
		return
	}

	boil.SetDB(db)
	booksG, err := models.Books().AllG(ctx)
	if err != nil {
		log.Printf("本の取得エラー: %v\n", err)
		return
	}

	booksWithCondition, err := models.Books(
		qm.Where("title LIKE ?", "%Go%"),
		qm.Limit(5),
	).All(ctx, db)
	if err != nil {
		log.Printf("クエリエラー: %v\n", err)
		return
	}

	// 結果の表示
	fmt.Println("=== 全ての本 ===")
	for _, b := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}

	fmt.Println("\n=== 条件付きクエリの結果 ===")
	for _, b := range booksWithCondition {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}

	// 検索関数の使用
	searchResult, err := searchBooks(ctx, db, "Programming")
	if err != nil {
		log.Printf("検索エラー: %v\n", err)
		return
	}

	fmt.Println("\n=== 検索結果 ===")
	for _, b := range searchResult {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", b.ID, b.Title, b.Author)
	}
	// ---------------------------

	// トランザクションの例
	// db.BeginTx: データベーストランザクションを開始するメソッド
	// ctx: コンテキスト（タイムアウトや中断の制御に使用）
	// nil: トランザクションのオプション（デフォルト設定を使用）
	// 戻り値のtxはトランザクションオブジェクト
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("トランザクション開始エラー: %v\n", err)
		return
	}
	// トランザクションのロールバックを保証
	defer tx.Rollback()

	// トランザクション内でのクエリ実行
	usersInTx, err := models.Users().All(ctx, tx)
	if err != nil {
		log.Printf("ユーザー取得エラー: %v\n", err)
		return
	}

	// トランザクションの利点：

	// 複数の操作をまとめて実行（全て成功するか、全て失敗するか）
	// データの一貫性を保証
	// エラー発生時に安全に元に戻せる
	// 使用例：

	// ユーザーの作成と同時に関連データも作成
	// 銀行取引（送金元と送金先の残高更新）
	// 在庫管理（注文処理と在庫数更新）

	// リレーションシップの例
	// 単一のユーザーを取得
	// 以下のようなSQLと同じ:
	// SELECT * FROM users LIMIT 1;
	user, err := models.Users().One(ctx, db)
	if err != nil {
		log.Printf("ユーザー取得エラー: %v\n", err)
		return
	}

	// ユーザーのお気に入り映画を取得
	movies, err := user.FavoriteMovies().All(ctx, db)
	if err != nil {
		log.Printf("お気に入り映画取得エラー: %v\n", err)
		return
	}

	// Eager loading（関連データの一括取得）の例
	// SQLでは以下のようなJOINクエリに相当:
	// SELECT users.*, movies.*
	// FROM users
	// JOIN movies ON users.id = movies.user_id;
	usersWithMovies, err := models.Users(qm.Load("FavoriteMovies")).All(ctx, db)
	if err != nil {
		log.Printf("ユーザーと映画の取得エラー: %v\n", err)
		return
	}
	// qmとは
	// import "github.com/volatiletech/sqlboiler/v4/queries/qm"
	// qmは"Query Modifier"（クエリ修飾子）の略
	// SQLクエリを構築するためのヘルパーパッケージ
	// SQLの各部分（WHERE, ORDER BY, LIMIT など）を簡単に書けるようにする
	// Load関数の基本
	// qm.Load("FavoriteMovies")
	// Loadは関連するデータを一緒に取得するための機能
	// 引数の"FavoriteMovies"は関連テーブルの名前

	// 結果の表示
	fmt.Println("\n=== ユーザーとお気に入り映画 ===")
	for _, u := range usersWithMovies {
		fmt.Printf("User ID: %d, Name: %s\n", u.ID, u.Name)
		if u.R != nil && u.R.FavoriteMovies != nil {
			fmt.Printf("お気に入り映画数: %d\n", len(u.R.FavoriteMovies))
			for _, m := range u.R.FavoriteMovies {
				fmt.Printf("  - Movie: %s\n", m.Title)
			}
		}
	}

	// トランザクションのコミット
	if err := tx.Commit(); err != nil {
		log.Printf("トランザクションコミットエラー: %v\n", err)
		return
	}
}
