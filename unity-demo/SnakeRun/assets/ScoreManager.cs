using UnityEngine;
using UnityEngine.SceneManagement;
using UnityEngine.UI; 

public class ScoreManager : MonoBehaviour
{
    public static int score = 0;
    public Text scoreText;

    private void Update()
    {
        // Update the text field with the current score
        scoreText.text = "Score:" + score.ToString();
    }


    //游戏重新开始，分数清零，重新加载游戏场景
    public void OnRestartButtonClick() // Connect this function to your button's onClick event in the inspector
    {
        SceneManager.LoadScene(SceneManager.GetActiveScene().buildIndex);
        score = 0;
    }
}
