using UnityEngine;

public class SnakeController : MonoBehaviour
{
    //每次跳动的高度
    public float jumpForce = 5f;
    private Rigidbody2D rb;
    public GameObject gameOverPanel;
    private bool isGameOver = false;

    public float upperLimit = 1000000000f; // Set this to the top of your screen
    public float lowerLimit = -1000000000f; // Set this to the bottom of your screen

    //public AudioSource jumpSFX;

    private void Start()
    {
        rb = GetComponent<Rigidbody2D>();
    }

    private void Update()
    {
        if (isGameOver) return;

        // Check if Snake is out of bounds【小蛇超出页面也触发游戏结束】
        // if (transform.position.y > 50f || transform.position.y < -50f)
        // {
        //     //Debug.Log(transform.position.y);
        //     GameOver();
        // }

        if (Input.GetKeyDown(KeyCode.Space))
        {
            Jump();
        }
    }

    private void Jump()
    {
        rb.linearVelocity = Vector2.up * jumpForce;
        //jumpSFX.Play();
    }


    //与Trigger部分碰撞时，触发分数加操作
    private void OnTriggerEnter2D(Collider2D collision)
    {
        // if (isGameOver) return; // Add this line

        //Debug.Log("Score: " + ScoreManager.score);
        ScoreManager.score++;
    }

    //当snake与Barrier相撞时，游戏结束
    private void OnCollisionEnter2D(Collision2D collision)
    {
        if (collision.gameObject.CompareTag("Barrier"))
        {
            // Game over
            GameOver();
        }
    }

    private void GameOver()
    {
        isGameOver = true; // Add this line

        // Freeze the Snake's motion
        rb.linearVelocity = Vector2.zero;

        if (gameOverPanel != null)
        {
            //展示游戏结束页面
            gameOverPanel.SetActive(true);
        }
    }
}
