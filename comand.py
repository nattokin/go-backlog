import re
import sys
from typing import List, Tuple

def simplify_delegation_comments(file_path: str):
    """
    Goファイルの委譲メソッドのコメントを簡略化し、結果を標準出力に出力します。

    対象のコメントパターン:
    // WithQueryActivityTypeIDs delegates to QueryOptionService.
    // WithFormArchived returns a form option that sets the `archived` field for the project.  <-- これも対象外だが、簡略化の余地があるためパターンを追加

    置換後のコメントパターン:
    // WithQueryActivityTypeIDs creates a query option to filter by activity type IDs.
    """
    
    # 処理対象のファイルを開く
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"Error: File not found at {file_path}", file=sys.stderr)
        return
    except Exception as e:
        print(f"Error reading file: {e}", file=sys.stderr)
        return

    # 1. 'delegates to XxxService.' パターンを検出・置換
    #   - 1: コメント行全体 (例: // WithQueryActivityTypeIDs delegates to QueryOptionService.)
    #   - 2: メソッド名 (例: WithQueryActivityTypeIDs)
    #   - 3: 委譲先のサービス名 (例: QueryOptionService)
    delegate_pattern = re.compile(
        r"^\s*//\s*(\w+)\s+delegates\s+to\s+([A-Za-z0-9]+)\.\s*$",
        re.MULTILINE
    )

    def replace_delegate_comment(match):
        method_name = match.group(1) # 例: WithQueryActivityTypeIDs
        service_name = match.group(2) # 例: QueryOptionService
        
        # 'With' 以降の部分を抽出（例: QueryActivityTypeIDs）
        func_part = method_name.removeprefix("With") 
        
        # 命名規則から単語を分割し、意味のある説明文を生成
        # 例: QueryActivityTypeIDs -> creates a query option for Activity Type IDs.
        # 例: FormPassword -> creates a form option for the password field.
        
        # 単語の境界を見つけるためのロジック (大文字の直前)
        parts = re.findall(r'[A-Z][a-z0-9]*', func_part)
        
        # 最初の単語 (Query/Form) に基づいて説明文のベースを決定
        if parts:
            option_type = parts[0].lower() # query or form
            
            # 残りの単語をスペース区切りに変換
            description = ' '.join(parts[1:]) 
            
            if option_type in ["query", "form"]:
                # 'option' を追加し、後の単語を小文字のまま連結
                return f"// {method_name} creates a {option_type} option for {description}."
            
        # 検出できなかった場合や、特殊なケースの場合は元のコメントをそのまま返す
        return match.group(0)

    # 置換を実行
    new_content = delegate_pattern.sub(replace_delegate_comment, content)
    
    # 2. 'returns a form option that sets the `key` field for the project.' などのパターンを簡略化 (ProjectOptionServiceのForm委譲メソッドなど)
    # 委譲元のコメントが簡潔な場合、そのまま残す
    
    # 以下のパターンも対象だが、前の置換で大部分が処理されるため、特に処理を追加しない。

    print(new_content)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python simplify_comments.py <path_to_go_file>")
        sys.exit(1)
        
    file_path = sys.argv[1]
    simplify_delegation_comments(file_path)

